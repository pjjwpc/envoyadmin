package models

import "time"

type SysMenu struct {
	ID         int        `gorm:"primary_key;auto_increment"`
	ParentID   int        `gorm:"not null"`
	TreePath   string     `gorm:"size:255"`
	Name       string     `gorm:"size:64;not null;default:''"`
	Type       int        `gorm:"not null"`
	Path       string     `gorm:"size:128;default:''"`
	Component  string     `gorm:"size:128"`
	Perm       string     `gorm:"size:128"`
	Visible    int        `gorm:"not null;default:1"`
	Sort       int        `gorm:"default:0"`
	Icon       string     `gorm:"size:64;default:''"`
	Redirect   string     `gorm:"size:128"`
	CreateTime *time.Time `gorm:"default:null"`
	UpdateTime *time.Time `gorm:"default:null"`
	AlwaysShow bool       `gorm:"always_show,omitempty"`
	KeepAlive  bool       `gorm:"keep_alive,omitempty"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
