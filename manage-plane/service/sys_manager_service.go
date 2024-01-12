package service

import (
	"errors"
	"log"
	"manage-plane/db"
	"manage-plane/models"
	"manage-plane/models/input"
	"manage-plane/utils"
)

type SysManagerService struct {
}

func (its *SysManagerService) Login(loginput *input.LoginInput) (token string, err error) {
	var dbUser models.SysUser
	err = db.Orm.Model(&models.SysUser{}).Where("username=?", loginput.UserName).Find(&dbUser).Error
	if err != nil {
		return "", errors.New("登录失败，用户不存在")
	}

	// 使用md5加密 ，进行密码比对
	mdspwd := utils.Md5V(loginput.Password)
	log.Println(loginput.Password, loginput.UserName, mdspwd, dbUser.Password)
	if mdspwd != dbUser.Password {
		return "", errors.New("密码错误")
	}
	// 查询用户角色
	var dbRole models.SysRole
	err = db.Orm.Raw(`SELECT sr.*
FROM sys_role sr INNER JOIN sys_user_role sur ON sr.id=sur.role_id
WHERE sur.user_id=?`, dbUser.ID).First(&dbRole).Error

	// 生成token
	token, err = utils.GenerateToken(dbRole.Code, dbUser.ID, dbUser.Nickname)
	if err != nil {
		return "", err
	}
	return token, nil

}
func (service *SysManagerService) GetUserInfo(userid int) (user models.SysUser, err error) {
	err = db.Orm.Model(&models.SysUser{}).Where("id=? ", userid).Preload("Roles").Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (service *SysManagerService) GetMenus(userid int, role string) (menus []models.Menu, err error) {
	// 根据用户id查询用户菜单
	var dbMenus []models.SysMenu
	if role == "ADMIN" || role == "ROOT" {
		err = db.Orm.Raw(`SELECT sm.* FROM sys_menu sm`).Scan(&dbMenus).Error
	} else {
		err = db.Orm.Raw(`SELECT sm.*
						  FROM sys_menu sm
						  INNER JOIN sys_role_menu srm ON srm.menu_id=sm.id
						  INNER JOIN  sys_user_role syr ON syr.role_id =srm.role_id
						  WHERE syr.user_id = ? `, userid).Scan(&dbMenus).Error
	}
	if err != nil {
		return nil, err
	}
	for _, dbMenu := range dbMenus {
		if dbMenu.ParentID == 0 {
			var menu models.Menu
			menu.Path = dbMenu.Path
			menu.Component = dbMenu.Component
			menu.Redirect = dbMenu.Redirect
			menu.Name = dbMenu.Name
			var meta models.Menumeta
			meta.AlwaysShow = true
			meta.Hidden = dbMenu.Visible == 0
			meta.Icon = dbMenu.Icon
			meta.KeepAlive = true
			meta.Title = dbMenu.Name
			meta.Roles = []string{role}
			menu.Children = []models.Menu{}
			for _, childMenu := range getChildMenus(dbMenu.ID, dbMenus, role) {
				menu.Children = append(menu.Children, childMenu)
			}
			menu.Meta = meta
			menus = append(menus, menu)
		}
	}
	return menus, nil
}

func getChildMenus(parentId int, menus []models.SysMenu, role string) []models.Menu {
	var childMenus []models.Menu
	for _, menu := range menus {
		if menu.ParentID == parentId {
			var childMenu models.Menu
			childMenu.Path = menu.Path
			childMenu.Component = menu.Component
			childMenu.Redirect = menu.Redirect
			childMenu.Name = menu.Name
			var meta models.Menumeta
			meta.AlwaysShow = true
			meta.Hidden = menu.Visible == 0
			meta.Icon = menu.Icon
			meta.KeepAlive = true
			meta.Title = menu.Name
			meta.Roles = []string{role}
			childMenu.Meta = meta
			childMenu.Children = getChildMenus(menu.ID, menus, role)
			if len(childMenu.Children) <= 0 {
				childMenu.Children = []models.Menu{}
			}
			childMenus = append(childMenus, childMenu)
		}
	}
	return childMenus
}
