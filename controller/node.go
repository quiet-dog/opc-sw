package controller

import (
	"sw/global"
	"sw/model/node"
	"sw/opc"

	"fmt"

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

type RecData struct {
	DeviceType              string                  `json:"deviceType"`
	DeviceId                string                  `json:"deviceId"`
	EnvironmentAlarmInfoDTO EnvironmentAlarmInfoDTO `json:"environmentAlarmInfoDTO"`
	EquipmentInfoDTO        EquipmentInfoDTO        `json:"equipmentInfoDTO"`
}

type EnvironmentAlarmInfoDTO struct {
	EnvironmentId    int     `json:"environmentId"`
	Value            float64 `json:"value"`
	Unit             string  `json:"unit"`
	Power            float64 `json:"power"`
	WaterValue       float64 `json:"waterValue"`
	ElectricityValue float64 `json:"electricityValue"`
}

type EquipmentInfoDTO struct {
	EquipmentId int     `json:"equipmentId"`
	ThresholdId int     `json:"thresholdId"`
	SensorName  string  `json:"sensorName"`
	Value       float64 `json:"value"`
}

func RecDataApi(c *gin.Context) {
	c.JSON(200, gin.H{"message": "数据发送成功"})
	var recData RecData
	if err := c.ShouldBindJSON(&recData); err != nil {
		// c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	reg := "deviceType-"
	var v float64
	if recData.EnvironmentAlarmInfoDTO.EnvironmentId != 0 {
		reg += "环境档案-environmentId-" + fmt.Sprint(recData.EnvironmentAlarmInfoDTO.EnvironmentId)
		v = recData.EnvironmentAlarmInfoDTO.Value
	}
	if recData.EquipmentInfoDTO.ThresholdId != 0 {
		reg += "设备档案-thresholdId-" + fmt.Sprint(recData.EquipmentInfoDTO.ThresholdId)
		v = recData.EquipmentInfoDTO.Value
	}
	var nodeModel node.NodeModel
	global.DB.Where("param like ?", "%"+reg+"%").First(&nodeModel)
	if nodeModel.ID == 0 {
		// c.JSON(400, gin.H{"error": "没有找到对应的节点"})
		return
	}

	send := global.RecHandler{}
	send.Type = global.DEVICEDATA
	var sd opc.Data
	sd.DataType = "FLOAT64"
	sd.ID = uint64(nodeModel.ID)
	sd.Value = v
	send.Data = sd
	global.RecChanel <- send
	// c.JSON(200, gin.H{"message": "数据发送成功"})
}

type SendThresholdDTO struct {
	Threshold       ThresholdEntity        `json:"threshold"`
	ThresholdValues []ThresholdValueEntity `json:"thresholdValues"`
}

type ThresholdEntity struct {
	ThresholdID    int64  `json:"threshold_id" gorm:"column:threshold_id;primaryKey;autoIncrement"`
	EquipmentID    int64  `json:"equipment_id" gorm:"column:equipment_id"`
	SensorName     string `json:"sensor_name" gorm:"column:sensor_name"`
	SensorModel    string `json:"sensor_model" gorm:"column:sensor_model"`
	EquipmentIndex string `json:"equipment_index" gorm:"column:equipment_index"`
	Unit           string `json:"unit" gorm:"column:unit"`
	Code           string `json:"code" gorm:"column:code"`
	PurchaseDate   string `json:"purchase_date" gorm:"column:purchase_date"`
	OutID          string `json:"out_id" gorm:"column:out_id"`
}

type ThresholdValueEntity struct {
	ThresholdID int64   `json:"threshold_id" gorm:"column:threshold_id;primaryKey;autoIncrement"`
	Min         float64 `json:"min" gorm:"column:min"`
	Max         float64 `json:"max" gorm:"column:max"`
	Level       string  `json:"level" gorm:"column:level"`
}

func RecYuZhiApi(c *gin.Context) {
	c.JSON(200, gin.H{"message": "数据发送成功"})
	var recData SendThresholdDTO
	if err := c.ShouldBindJSON(&recData); err != nil {
		return
	}
	var send global.RecHandler
	send.Type = global.YUZHI
	send.Data = recData
	global.RecChanel <- send

}

func RecBaoJingApi(c *gin.Context) {

}
