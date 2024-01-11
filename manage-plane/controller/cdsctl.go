package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClusterController struct {
}

func (its *ClusterController) Get(c *gin.Context) {
}

func ErrMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": "-1",
	})
}

func OkMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": "200",
	})
}

func (its *ClusterController) Post(c *gin.Context) {
}

func (its *ClusterController) Patch(c *gin.Context) {
}

func (its *ClusterController) Delete(c *gin.Context) {
}

func (its *ClusterController) GetVersionByEnvoyClusterId(c *gin.Context) {
}

func (its *ClusterController) GetUnUseCds(c *gin.Context) {
}
