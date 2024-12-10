package node

import (
	"fmt"
	"sw/global"

	"gorm.io/gorm"
)

type NodeModel struct {
	gorm.Model
	NodeId    string `json:"nodeId"`
	Param     string `json:"param"`
	ServiceId uint   `json:"serviceId"`
}

func (n *NodeModel) AfterCreate(tx *gorm.DB) error {
	err := global.OpcGateway.AddNode(fmt.Sprintf("%d", n.ServiceId), n.NodeId)
	return err
}

type AddNode struct {
	NodeId    string `json:"nodeId"`
	Param     string `json:"param"`
	ServiceId uint   `json:"serviceId"`
}

type UpdateNode struct {
	Id uint `json:"id"`
	AddNode
}

func LoadAddNode(add AddNode) *NodeModel {
	return &NodeModel{
		NodeId:    add.NodeId,
		Param:     add.Param,
		ServiceId: add.ServiceId,
	}
}

func LoadUpdateNode(update UpdateNode) *NodeModel {
	var n NodeModel
	global.DB.First(&n, update.Id)
	n.NodeId = update.NodeId
	n.Param = update.Param
	n.ServiceId = update.ServiceId
	return &n
}

func (n *NodeModel) Create() {
	global.DB.Create(n)
}

func (n *NodeModel) Update() {
	global.DB.Save(n)
}

func (n *NodeModel) Delete() {
	global.DB.Delete(n)
}
