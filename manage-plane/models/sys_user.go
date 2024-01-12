package models

import "time"

type SysUser struct {
	ID         int        `gorm:"primary_key;auto_increment"`
	Username   string     `gorm:"size:64"`
	Nickname   string     `gorm:"size:64"`
	Gender     int        `gorm:"default:1"`
	Password   string     `gorm:"size:100"`
	DeptID     int        `gorm:"default:null"`
	Avatar     string     `gorm:"size:255;default:''"`
	Mobile     string     `gorm:"size:20"`
	Status     int        `gorm:"default:1"`
	Email      string     `gorm:"size:128"`
	Deleted    int        `gorm:"default:0"`
	CreateTime *time.Time `gorm:"default:null"`
	UpdateTime *time.Time `gorm:"default:null"`
	Roles      []SysRole  `gorm:"many2many:sys_user_role;foreignKey:id;joinForeignKey:user_id;joinReferences:role_id"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
