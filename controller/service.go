package controller

import (
	"os/exec"
	"sw/global"
	"sw/model/service"

	"github.com/gin-gonic/gin"
)

func CreateService(c *gin.Context) {
	var cService service.AddService
	if err := c.ShouldBindJSON(&cService); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	serviceModel := service.LoadAddService(cService)
	serviceModel.Create()
	c.JSON(200, serviceModel)
}

func UpdateService(c *gin.Context) {
	var uService service.UpdateService
	if err := c.ShouldBindJSON(&uService); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	serviceModel := service.LoadUpdateService(uService)
	serviceModel.Update()
	c.JSON(200, serviceModel)
}

func DeleteService(c *gin.Context) {
	var dService service.ServiceModel
	if err := c.ShouldBindJSON(&dService); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	dService.Delete()
	c.JSON(200, dService)
}

func GetServiceList(c *gin.Context) {
	var services []*service.ServiceModel
	global.DB.Find(&services)
	c.JSON(200, services)
}

func RestSys(c *gin.Context) {
	// 调用命令行  systemctl restart opc

	// 是linux
	cmd := exec.Command("systemctl", "restart", "opc")
	err := cmd.Run()
	if err != nil {
		c.JSON(500, gin.H{"error": "重启服务失败"})
		return
	}
	c.JSON(200, gin.H{"message": "重启服务成功"})

}
