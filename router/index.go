package router

import (
	"sw/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/node", controller.GetNodeList)
	r.POST("/node", controller.CreateNode)
	r.POST("/node/update", controller.UpdateNode)
	r.POST("/node/delete", controller.DeleteNode)
	r.GET("/service", controller.GetServiceList)
	r.POST("/service", controller.CreateService)
	r.POST("/service/update", controller.UpdateService)
	r.POST("/service/delete", controller.DeleteService)
	return r
}
