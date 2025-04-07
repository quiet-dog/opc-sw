package node

import (
	"encoding/json"
	"fmt"
	"sw/global"
	"sw/opc"
	"time"

	"gorm.io/gorm"
)

/*
*
const deviceType = ref("")

	const deviceName = ref("")
	const bmsDeviceName = ref("")
	const bmsArea =ref("")
	const bmsLabel = ref("")
	const emsAare = ref("")
*/
type NodeModel struct {
	gorm.Model
	NodeId        string      `json:"nodeId"`
	Param         string      `json:"param"`
	ServiceId     uint        `json:"serviceId"`
	Description   string      `json:"description"`
	Time          time.Time   `gorm:"-" json:"time"`
	Value         interface{} `gorm:"-" json:"value"`
	Type          string      `gorm:"-" json:"type"`
	Extend        string      `json:"extend"`
	DeviceName    string      `json:"deviceName"`
	Key           string      `json:"key"`
	DeviceType    string      `json:"deviceType"`
	BmsDeviceName string      `json:"bmsDeviceName"`
	BmsArea       string      `json:"bmsArea"`
	BmsLabel      string      `json:"bmsLabel"`
	EmsAare       string      `json:"emsAare"`
}

func (n *NodeModel) AfterCreate(tx *gorm.DB) error {
	// err := global.OpcGateway.AddNode(fmt.Sprintf("%d", n.ServiceId), opc.NodeId{
	// 	ID:   uint64(n.ID),
	// 	Node: n.NodeId,
	// })
	return nil
}

func (n *NodeModel) AfterFind(tx *gorm.DB) error {
	var notify opc.Data
	b, err := global.Redis.Get(global.Ctx, fmt.Sprintf("%d", n.ID)).Result()
	if err != nil {
		fmt.Println("获取redis错误", err.Error())
		return nil
	}
	err = json.Unmarshal([]byte(b), &notify)
	if err != nil {
		fmt.Println("json unmarshal error", err.Error())
		return nil
	}
	n.Time = notify.SourceTime
	n.Value = notify.Value
	n.Type = notify.DataType
	return nil
}

type AddNode struct {
	NodeId        string `json:"nodeId"`
	Param         string `json:"param"`
	ServiceId     uint   `json:"serviceId"`
	Description   string `json:"description"`
	Extend        string `json:"extend"`
	DeviceName    string `json:"deviceName"`
	Key           string `json:"key"`
	DeviceType    string `json:"deviceType"`
	BmsDeviceName string `json:"bmsDeviceName"`
	BmsArea       string `json:"bmsArea"`
	BmsLabel      string `json:"bmsLabel"`
	EmsAare       string `json:"emsAare"`
}

type UpdateNode struct {
	Id uint `json:"id"`
	AddNode
}

func LoadAddNode(add AddNode) *NodeModel {
	return &NodeModel{
		NodeId:        add.NodeId,
		Param:         add.Param,
		ServiceId:     add.ServiceId,
		Description:   add.Description,
		Extend:        add.Extend,
		DeviceName:    add.DeviceName,
		DeviceType:    add.DeviceType,
		Key:           add.Key,
		BmsDeviceName: add.BmsDeviceName,
		BmsArea:       add.BmsArea,
		BmsLabel:      add.BmsLabel,
		EmsAare:       add.EmsAare,
	}
}

func LoadUpdateNode(update UpdateNode) *NodeModel {
	var n NodeModel
	global.DB.First(&n, update.Id)
	n.NodeId = update.NodeId
	n.Param = update.Param
	n.ServiceId = update.ServiceId
	n.Description = update.Description
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
