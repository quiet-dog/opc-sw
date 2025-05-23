package router

import (
	"sw/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	a := r.Group("/api")
	r.GET("/ws", controller.Connect)
	a.POST("/node/list", controller.GetNodeList)
	a.POST("/node", controller.CreateNode)
	a.POST("/node/update", controller.UpdateNode)
	a.POST("/node/delete", controller.DeleteNode)
	a.GET("/service", controller.GetServiceList)
	a.POST("/service", controller.CreateService)
	a.POST("/service/update", controller.UpdateService)
	a.POST("/service/delete", controller.DeleteService)

	a.POST("/recDataApi", controller.RecDataApi)
	a.POST("/recYuZhiApi", controller.RecYuZhiApi)
	return r
}
