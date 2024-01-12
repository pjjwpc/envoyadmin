package models

import "time"

type SysRole struct {
	ID         int64      `gorm:"primary_key;auto_increment"`
	Name       string     `gorm:"size:64;not null;default:''"`
	Code       string     `gorm:"size:32"`
	Sort       int        `gorm:"default:null"`
	Status     int        `gorm:"default:1"`
	DataScope  int        `gorm:"default:null"`
	Deleted    int        `gorm:"not null;default:0"`
	CreateTime *time.Time `gorm:"default:null"`
	UpdateTime *time.Time `gorm:"default:null"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

type SysRoleMenu struct {
	RoleID int64 `gorm:"primary_key"`
	MenuID int64 `gorm:"primary_key"`
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}

type SysUserRole struct {
	UserID int64 `gorm:"primary_key"`
	RoleID int64 `gorm:"primary_key"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}
