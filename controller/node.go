package controller

import (
	"sw/global"
	"sw/model/node"
	"sw/opc"

	"fmt"
	"time"

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
	EnvironmentAlarmInfoDTO EnvironmentAlarmInfoDTO `json:"environmentAlarmInfo"`
	EquipmentInfoDTO        EquipmentInfoDTO        `json:"equipmentInfo"`
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
	// 打印body数据

	var recData RecData
	if err := c.ShouldBindJSON(&recData); err != nil {
		// c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(recData)
	reg := "deviceType-"
	fmt.Println("======", recData.EnvironmentAlarmInfoDTO.EnvironmentId)
	var v float64
	if recData.EnvironmentAlarmInfoDTO.EnvironmentId != 0 {
		reg += fmt.Sprintf("环境档案-environmentId-%d", recData.EnvironmentAlarmInfoDTO.EnvironmentId)
		v = recData.EnvironmentAlarmInfoDTO.Value
	}
	if recData.EquipmentInfoDTO.ThresholdId != 0 {
		reg += fmt.Sprintf("设备档案-thresholdId-%d", recData.EquipmentInfoDTO.ThresholdId)
		v = recData.EquipmentInfoDTO.Value
	}
	fmt.Println("reg:", reg)
	var nodeModel node.NodeModel
	global.DB.Where("param like ?", "%"+reg+"%").First(&nodeModel)
	if nodeModel.ID == 0 {
		fmt.Println("没有找到数据1")
		// c.JSON(400, gin.H{"error": "没有找到对应的节点"})
		return
	}

	send := global.RecHandler{}
	send.Type = global.DEVICEDATA
	var sd opc.Data
	sd.DataType = "FLOAT64"
	sd.ID = uint64(nodeModel.ID)
	sd.Value = v
	sd.SourceTime = time.Now()
	sd.Param = nodeModel.Param
	send.Data = sd
	fmt.Println("发送数据1")
	global.RecChanel <- send
	global.MoniChannel <- sd
	// c.JSON(200, gin.H{"message": "数据发送成功"})
}

type SendThresholdDTO struct {
	Threshold       ThresholdEntity        `json:"threshold"`
	ThresholdValues []ThresholdValueEntity `json:"thresholdValues"`
	// Node            node.NodeModel         `json:"node"`
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

type SendEnvrionmentDTO struct {
	Environment EnvironmentEntity        `json:"environment"`
	Values      []AlarmlevelDetailEntity `json:"values"`
}

type EnvironmentEntity struct {
	EnvironmentID   int64  `json:"environmentId" gorm:"column:environment_id;primaryKey;autoIncrement"` // 主键 ID
	Description     string `json:"description" gorm:"column:description"`                               // 描述
	MonitoringPoint string `json:"monitoringPoint" gorm:"column:monitoring_point"`                      // 监测点位
	// EnvironmentIndex string  `json:"environment_index" gorm:"column:environment_index"`                // 环境指标（已注释）
	Tag        string  `json:"tag" gorm:"column:tag"`               // 位号
	Type       string  `json:"type" gorm:"column:type"`             // 类型
	Scope      string  `json:"scope" gorm:"column:scope"`           // 范围
	ESignal    string  `json:"eSignal" gorm:"column:eSignal"`       // 信号
	Supplier   string  `json:"supplier" gorm:"column:supplier"`     // 设备仪表供应商
	Model      string  `json:"model" gorm:"column:model"`           // 设备仪表型号
	PLCAddress string  `json:"plcAddress" gorm:"column:plcAddress"` // PLC地址
	EArea      string  `json:"eArea" gorm:"column:e_area"`          // 区域
	Value      float64 `json:"value" gorm:"column:value"`           // 数值
	Unit       string  `json:"unit" gorm:"column:unit"`             // 单位
	UnitName   string  `json:"unitName" gorm:"column:unit_name"`    // 监测指标
}

type AlarmlevelDetailEntity struct {
	AlarmlevelID  int64   `json:"alarmlevelId" gorm:"column:alarmlevel_id;primaryKey;autoIncrement"` // 主键 ID
	EnvironmentID int64   `json:"environmentId" gorm:"column:environment_id"`                        // 环境 ID
	Min           float64 `json:"min" gorm:"column:min"`                                             // 最小值
	Max           float64 `json:"max" gorm:"column:max"`                                             // 最大值
	Level         string  `json:"level" gorm:"column:level"`                                         // 级别
	Unit          string  `json:"unit" gorm:"column:unit"`                                           // 单位
}

func RecYuZhiApi(c *gin.Context) {
	c.JSON(200, gin.H{"message": "数据发送成功"})
	var recData SendThresholdDTO
	if err := c.ShouldBindJSON(&recData); err != nil {
		return
	}
	var send global.RecHandler
	send.Type = global.YUZHI
	var nodeModle node.NodeModel
	global.DB.Where("param like ?", "%设备档案-thresholdId-"+fmt.Sprint(recData.Threshold.ThresholdID)+"%").First(&nodeModle)
	send.Data = map[string]interface{}{
		"threshold":       recData.Threshold,
		"thresholdValues": recData.ThresholdValues,
		"node":            nodeModle,
	}

	global.RecChanel <- send

}

func RecEnvYuZhiApi(c *gin.Context) {
	c.JSON(200, gin.H{"message": "数据发送成功"})
	var recData SendEnvrionmentDTO
	if err := c.ShouldBindJSON(&recData); err != nil {
		return
	}
	var send global.RecHandler
	send.Type = global.YUZHI
	var nodeModle node.NodeModel
	global.DB.Where("param like ?", "%环境档案-environmentId-"+fmt.Sprint(recData.Environment.EnvironmentID)+"%").First(&nodeModle)
	send.Data = map[string]interface{}{
		"environment":       recData.Environment,
		"environmentValues": recData.Values,
		"node":              nodeModle,
	}

	global.RecChanel <- send

}

type AlarmEvent struct {
	EventID          int         `json:"eventId"`
	Type             string      `json:"type"`
	EquipmentID      int         `json:"equipmentId"`
	EquipmentValue   float64     `json:"equipmentValue"`
	Equipment        interface{} `json:"equipment"`
	MaterialsID      int         `json:"materialsId"`
	MaterialsValue   float64     `json:"materialsValue"`
	Materials        interface{} `json:"materials"`
	EnvironmentID    int         `json:"environmentId"`
	Environment      interface{} `json:"environment"`
	EnvironmentValue float64     `json:"environmentValue"`
	MaterialsValueID int         `json:"materialsValueId"`
	Level            string      `json:"level"`
	Description      string      `json:"description"`
	HandlerID        int         `json:"handlerId"`
	Handler          interface{} `json:"handler"`
	AlarmLevelID     int         `json:"alarmlevelId"`
	ThresholdID      int         `json:"thresholdId"`
	CraftNodeID      int         `json:"craftNodeId"`
	CraftNode        interface{} `json:"craftNode"`
	Threshold        interface{} `json:"threshold"`
	Emergencys       interface{} `json:"emergencys"`
	SOPs             interface{} `json:"sops"`
	CreateTime       string      `json:"createTime"`
}

func RecBaoJingApi(c *gin.Context) {
	c.JSON(200, gin.H{"message": "数据发送成功"})
	// 打印body数据

	var recData AlarmEvent
	if err := c.ShouldBindJSON(&recData); err != nil {
		return
	}

	var send global.RecHandler
	send.Type = global.BAOJING
	var nodeModle node.NodeModel
	if recData.Type == "设备报警" {
		global.DB.Where("param like ?", "%设备档案-thresholdId-"+fmt.Sprint(recData.ThresholdID)+"%").First(&nodeModle)
	} else {
		global.DB.Where("param like ?", "%环境档案-environmentId-"+fmt.Sprint(recData.EnvironmentID)+"%").First(&nodeModle)
	}

	send.Data = map[string]interface{}{
		"event": recData,
		"node":  nodeModle,
	}

	global.RecChanel <- send
}

// 控制动画
func Animation(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Received animation control data:", data)
	var send global.RecHandler
	send.Type = global.ANIMATION
	send.Data = data

	global.RecChanel <- send
	c.JSON(200, gin.H{"message": "动画控制数据发送成功"})
}

func KetiSan(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var send global.RecHandler
	send.Type = global.KETISAN
	send.Data = data
	global.RecChanel <- send
	c.JSON(200, gin.H{"message": "数据发送成功"})
}
