package router

import (
	"net/http"

	api "manage-plane/controller"
	"manage-plane/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	route := gin.Default()

	route.POST("/login", api.ApiGroupApp.SysManagerApi.Login)
	systemManager := route.Group("/sys")
	systemManager.Use(middleware.AuthMiddleward())
	{
		systemManager.GET("/UserInfo", api.ApiGroupApp.SysManagerApi.UserInfo)
		systemManager.GET("/GetMenus", api.ApiGroupApp.SysManagerApi.GetMenus)
		systemManager.DELETE("/logout", api.ApiGroupApp.SysManagerApi.Logout)
	}
	envoyCluster := route.Group("/envoy/envoycluster")
	envoyCluster.Use(middleware.AuthMiddleward())
	{
		envoyCluster.GET("/", api.ApiGroupApp.EnvoyClusterApi.Get)
		envoyCluster.POST("/", api.ApiGroupApp.EnvoyClusterApi.Post)
		envoyCluster.DELETE("/", api.ApiGroupApp.EnvoyClusterApi.Delete)
		envoyCluster.PATCH("/", api.ApiGroupApp.EnvoyClusterApi.Patch)
		envoyCluster.GET("/GetNode", api.ApiGroupApp.EnvoyClusterApi.GetNode)
	}
	envoyNode := route.Group("/envoy/envoynode")
	envoyNode.Use(middleware.AuthMiddleward())
	{
		envoyNode.GET("/", api.ApiGroupApp.EnvoyNodeApi.Get)
		envoyNode.POST("/", api.ApiGroupApp.EnvoyNodeApi.Post)
		envoyNode.PATCH("/", api.ApiGroupApp.EnvoyNodeApi.Patch)
		envoyNode.DELETE("/disable", api.ApiGroupApp.EnvoyNodeApi.DisableNode)
		envoyNode.DELETE("/", api.ApiGroupApp.EnvoyNodeApi.RemoveNode)
	}
	cluster := route.Group("/envoy/cds")
	cluster.Use(middleware.AuthMiddleward())
	{
		cluster.GET("/", api.ApiGroupApp.CdsApi.Get)
		cluster.POST("/", api.ApiGroupApp.CdsApi.Post)
		cluster.PATCH("/", api.ApiGroupApp.CdsApi.Patch)
		cluster.DELETE("/", api.ApiGroupApp.CdsApi.Delete)
	}

	endpoint := route.Group("/envoy/eds")
	endpoint.Use(middleware.AuthMiddleward())
	{
		endpoint.GET("/", api.ApiGroupApp.EdsApi.Get)
		endpoint.POST("/", api.ApiGroupApp.EdsApi.Post)
		endpoint.PATCH("/", api.ApiGroupApp.EdsApi.Patch)
		endpoint.DELETE("/", api.ApiGroupApp.EdsApi.Delete)
	}

	route.GET("/status", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "200",
		})
	})

	return route
}
