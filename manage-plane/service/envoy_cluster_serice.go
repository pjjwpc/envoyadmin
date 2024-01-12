package service

import (
	"log"
	"manage-plane/db"
	"manage-plane/models"
)

type EnvoyClusterService struct {
}

func (ist *EnvoyClusterService) Count(id int64) int64 {
	tx := db.Orm.Model(&models.EnvoyClusterModel{})
	var total int64
	err := tx.Where("id=?", id).Count(&total).Error
	if err != nil {
		log.Println(err)
	}
	return total
}
