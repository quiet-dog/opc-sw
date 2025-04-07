package controller

import (
	"sw/global"
	"sw/model/node"

	"github.com/gin-gonic/gin"
)

func CreateNode(c *gin.Context) {
	var cNode node.AddNode
	if err := c.ShouldBindJSON(&cNode); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	nodeModel := node.LoadAddNode(cNode)
	nodeModel.Create()
	c.JSON(200, nodeModel)
}

func UpdateNode(c *gin.Context) {
	var uNode node.UpdateNode
	if err := c.ShouldBindJSON(&uNode); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	nodeModel := node.LoadUpdateNode(uNode)
	nodeModel.Update()
	c.JSON(200, nodeModel)
}

func DeleteNode(c *gin.Context) {
	var dNode node.NodeModel
	if err := c.ShouldBindJSON(&dNode); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	dNode.Delete()
	c.JSON(200, dNode)
}

type FindNodeParam struct {
	ServiceId uint `json:"serviceId"`
}

func GetNodeList(c *gin.Context) {
	var f FindNodeParam
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var nodes []*node.NodeModel

	if f.ServiceId == 0 {
		global.DB.Find(&nodes)
	} else {
		global.DB.Where("service_id = ?", f.ServiceId).Find(&nodes)
	}

	c.JSON(200, nodes)
}
