package service

import (
	"errors"
	"manage-plane/config"
	db "manage-plane/db"
	"manage-plane/models"
	"manage-plane/models/query"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
)

type ClusterService struct {
}

var envoyClusterService = EnvoyClusterService{}

func (its *ClusterService) Get(req query.CdsQuery) (list []models.Cds, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.PageIndex - 1)
	db := db.Orm.Model(&models.Cds{}).Where("1=1")

	if req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
	}

	if req.EnvoyClusterId > 0 {
		db = db.Where("envoy_cluster_id=? ", req.EnvoyClusterId)
	}
	err = db.Count(&total).Error
	if err != nil || total <= 0 {
		return list, total, err
	}
	err = db.Limit(limit).Offset(offset).Preload("EnvoyCluster").Order("envoy_cluster_id asc").Find(&list).Error

	return
}

func (its *ClusterService) Add(envoyClusterId int64, version string, uname string, cluster *cluster.Cluster) error {
	tx := db.Orm.Model(&models.Cds{})
	cstr, _ := config.ClusterFormat.Marshal(cluster)
	dbModel := models.Cds{}
	total := envoyClusterService.Count(envoyClusterId)
	if total <= 0 {
		return errors.New("envoy集群不存在，请检查")
	}
	err := tx.Where("envoy_cluster_id=? and name = ?", envoyClusterId, cluster.Name).
		Count(&total).Error
	if err != nil || total > 0 {
		return errors.New("已存在该Cluster")
	}

	dbModel.Name = cluster.Name
	dbModel.EnvoyClusterId = envoyClusterId
	dbModel.Type = cluster.GetType().String()
	dbModel.DnsLookupFamily = int(cluster.DnsLookupFamily)
	dbModel.LbPolicy = int(cluster.LbPolicy)
	dbModel.IsDelete = false
	dbModel.Enable = true
	dbModel.Version = version
	if len(cluster.HealthChecks) > 0 {
		dbModel.HealthCheck = true
	}
	dbModel.CreateTime = time.Now()
	dbModel.CreateUser = uname
	dbModel.UpdateTime = time.Now()
	dbModel.UpdateUser = uname
	dbModel.ValueData = string(cstr)
	txt := db.Orm.Begin()
	err = txt.Model(&models.Cds{}).Create(&dbModel).Error
	if err != nil {
		txt.Rollback()
		return err
	}
	// 这里需要再测试一下，怎样给每个cds设置单独的version
	// 现在的机制是所有的cluster使用统一的版本
	// 所以每次新增修改时都要把集群内所有的cluster版本更新一下

	txt.Commit()
	// todo 通知envoy更新
	return err
}
