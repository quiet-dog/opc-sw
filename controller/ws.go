package controller

import (
	"sw/global"

	"github.com/gin-gonic/gin"
)

func Connect(c *gin.Context) {
	socket, err := global.Upgrader.Upgrade(c.Writer, c.Request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 获取请求的参数url
	// 将url和socket绑定
	go func() {
		socket.ReadLoop()
	}()
}
