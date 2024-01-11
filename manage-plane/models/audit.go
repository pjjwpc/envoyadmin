package models

import (
	"time"
)

type Audit struct {
	CreateTime time.Time `gorm:"create_time" json:"createTime"`
	CreateUser string    `gorm:"create_user" json:"createUser"`
	UpdateTime time.Time `gorm:"update_time" json:"updateTime"`
	UpdateUser string    `gorm:"update_user" json:"updateUser"`
	IsDelete   bool      `gorm:"is_delete" json:"isDelete"`
}
