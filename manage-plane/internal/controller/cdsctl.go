package controller

import (
	"io"
	"log"
	"manage-plane/internal/config"
	"manage-plane/internal/service"
	"strconv"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/gin-gonic/gin"
)

type ClusterController struct {
	baseController
}

var clusterService = service.ClusterService{}

func (its *ClusterController) Get(c *gin.Context) {
}

func (its *ClusterController) Post(c *gin.Context) {
	u, r, _ := its.getUserInfo(c)
	if r != "OP" && r != "ADMIN" && r != "ROOT" {
		its.errMsg(c, "没有权限")
		return
	}
	clusterInput := cluster.Cluster{}
	envoy_cluster_id, err := strconv.ParseInt(c.Query("envoyClusterId"), 0, 64)
	version := c.Query("version")
	if err != nil || envoy_cluster_id == 0 {
		its.errMsg(c, "envoy集群参数未传，请检查")
		return
	}
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		its.errMsg(c, "参数有误，请检查")
		return
	}
	err = config.ClusterParser.Unmarshal(data, &clusterInput)
	if err != nil {
		log.Println(err)
		its.errMsg(c, "参数传递有误："+err.Error())
		return
	}
	err = clusterInput.ConnectTimeout.CheckValid()
	if err != nil {
		log.Println(err)
		its.errMsg(c, "connecttimeout配置错误")
		return
	}
	if clusterInput.GetType() == cluster.Cluster_STATIC {
		err = clusterInput.LoadAssignment.Validate()
		if err != nil {
			its.errMsg(c, "loadassignment配置错误，请检查")
			return
		}
	} else if clusterInput.GetType() == cluster.Cluster_EDS {
		err = clusterInput.EdsClusterConfig.Validate()
		if err != nil {
			its.errMsg(c, "eds配置有错误，请检查")
			return
		}
	}
	err = clusterService.Add(envoy_cluster_id, version, u, &clusterInput)

	if err != nil {
		log.Println(err)
		its.errMsg(c, err.Error())
		return
	}
	its.okMsg(c, "添加成功")

}

func (its *ClusterController) Patch(c *gin.Context) {
}

func (its *ClusterController) Delete(c *gin.Context) {
}

func (its *ClusterController) GetVersionByEnvoyClusterId(c *gin.Context) {
}

func (its *ClusterController) GetUnUseCds(c *gin.Context) {
}
