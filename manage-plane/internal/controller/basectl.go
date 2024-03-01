package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type baseController struct {
}

func (its *baseController) errMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": "-1",
	})
}

func (its *baseController) okMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": "200",
	})
}
func (its *baseController) getUserInfo(c *gin.Context) (username, role string, userId int) {
	userAny, exits := c.Get("userId")
	if !exits {
		log.Println("获取用户id失败")
		return
	}
	// 将userId 装换为int类型
	userIdStr := userAny.(float64)
	userId = int(userIdStr)
	usernameAny, exits := c.Get("username")
	if !exits {
		log.Println("获取用户名失败")
	}
	username = usernameAny.(string)
	roleAny, exits := c.Get("role")
	if !exits {
		log.Println("获取用户角色失败")
	}
	role = roleAny.(string)
	return
}
