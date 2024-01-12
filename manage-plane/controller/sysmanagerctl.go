package controller

import (
	"log"
	"manage-plane/models"
	"manage-plane/models/input"
	"manage-plane/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SysManagerController struct{}

var sysService = service.SysManagerService{}

func (its *SysManagerController) Login(c *gin.Context) {
	var loginInput input.LoginInput
	c.ShouldBind(&loginInput)
	token, err := sysService.Login(&loginInput)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"data": "用户名或密码错误",
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"tokenType":   "Bearer",
			"accessToken": token,
		},
		"msg": "一切OK",
	})
}

func (its *SysManagerController) UserInfo(c *gin.Context) {
	userId, exits := c.Get("userId")
	if !exits {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "重新登陆",
		})
		return
	}
	// 将userId 装换为int类型
	userIdStr := userId.(float64)
	userIdi := int(userIdStr)
	user, err := sysService.GetUserInfo(userIdi)
	if err != nil {
		log.Println("用户信息获取失败", err)
	}
	var roleCode []string
	for _, role := range user.Roles {
		roleCode = append(roleCode, role.Code)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"userId":   userIdi,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"roles":    roleCode,
			"perms":    []string{"sys:user:edit", "sys:user:delete", "sys:user:add"},
		},
		"msg": "一切OK",
	})
}

func (its *SysManagerController) GetMenus(c *gin.Context) {
	userId, exits := c.Get("userId")
	if !exits {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "重新登陆",
		})
		return
	}
	role, exits := c.Get("role")
	if !exits {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "重新登陆",
		})
		return
	}
	// 将userId 装换为int类型
	userIdStr := userId.(float64)
	userIdi := int(userIdStr)
	roleStr := role.(string)

	menus, err := sysService.GetMenus(userIdi, roleStr)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{
			"code": "500",
			"data": []models.Menu{},
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "200",
		"data": menus,
		"msg":  "一切OK",
	})
}

func (its *SysManagerController) AddUser(c *gin.Context) {

}

func (its *SysManagerController) Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": "200",
		"msg":  "注销成功",
	})
}
