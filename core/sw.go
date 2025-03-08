package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sw/global"
	"sw/model/node"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/gorilla/websocket"
)

// WebSocketReadWriteCloser 是 gorilla/websocket.Conn 的适配器
type WebSocketReadWriteCloser struct {
	Conn *websocket.Conn
}

// Read 实现 io.Reader 接口
func (w *WebSocketReadWriteCloser) Read(p []byte) (int, error) {
	_, message, err := w.Conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	n := copy(p, message)
	return n, nil
}

// Write 实现 io.Writer 接口
func (w *WebSocketReadWriteCloser) Write(p []byte) (int, error) {
	err := w.Conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Close 实现 io.Closer 接口
func (w *WebSocketReadWriteCloser) Close() error {
	return w.Conn.Close()
}

// DeviceDTO represents a device with various properties
type DeviceDTO struct {
	DeviceType           string                  `json:"deviceType"`           // 设备类型
	DeviceID             int64                   `json:"deviceId"`             // 设备ID
	EnvironmentAlarmInfo EnvironmentAlarmInfoDTO `json:"environmentAlarmInfo"` // 环境档案数据信息
	EquipmentInfo        EquipmentInfoDTO        `json:"equipmentInfo"`        // 设备信息
}

// EnvironmentAlarmInfoDTO represents environment alarm information
type EnvironmentAlarmInfoDTO struct {
	EnvironmentID    int64   `json:"environmentId"`    // 设备ID
	Value            float64 `json:"value"`            // 数据
	Unit             string  `json:"unit"`             // 单位
	Power            float64 `json:"power"`            // 功耗
	WaterValue       float64 `json:"waterValue"`       // 用水量
	ElectricityValue float64 `json:"electricityValue"` // 用电量
}

// EquipmentInfoDTO represents equipment information

// EquipmentInfoDTO represents equipment information
type EquipmentInfoDTO struct {
	EquipmentID int64   `json:"equipmentId"` // 设备ID
	ThresholdID int64   `json:"thresholdId"` // 阈值传感器ID
	SensorName  string  `json:"sensorName"`  // 传感器名称
	Value       float64 `json:"value"`       // 传感器值
}

func InitSw() {
	// WebSocket 连接信息
	url := fmt.Sprintf("ws://%s:%s/ws", global.Config.Sw.Host, global.Config.Sw.Port)
	header := http.Header{}
	ctx := context.Background()

	go func() {
		// 使用 gorilla/websocket 连接到服务器
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{

					conn, _, err := websocket.DefaultDialer.Dial(url, header)
					if err != nil {
						fmt.Println("Failed to connect to WebSocket server: %v", err)
						time.Sleep(5 * time.Second)
						continue
					}

					// 包装 WebSocket 连接为 io.ReadWriteCloser
					rwc := &WebSocketReadWriteCloser{Conn: conn}

					// 使用 STOMP 客户端连接
					stompConn, err := stomp.Connect(rwc)
					if err != nil {
						log.Fatalf("Failed to connect to STOMP: %v", err)
					}
					defer stompConn.Disconnect()

					log.Println("Connected to STOMP server")

					c := global.OpcGateway.SubscribeOpc()
					for {
						select {
						case msg, ok := <-c:
							{
								if !ok {
									return
								}

								result := DeviceDTO{}
								nodeModel := node.NodeModel{}
								global.DB.Where("id = ?", msg.ID).First(&nodeModel)
								// 字符串切割
								r := strings.Split(nodeModel.Param, "-")
								// deviceType := ""
								// environmentId := 0
								// thresholdId := 0
								// equipmentId := 0
								if len(r) >= 4 {
									for i := 0; i < len(r); i++ {
										if r[i] == "deviceType" {
											result.DeviceType = r[i+1]
										} else if r[i] == "environmentId" {
											if id, err := strconv.Atoi(r[i+1]); err == nil {
												result.EnvironmentAlarmInfo.EnvironmentID = int64(id)
											}

										} else if r[i] == "thresholdId" {
											if id, err := strconv.Atoi(r[i+1]); err == nil {
												result.EquipmentInfo.ThresholdID = int64(id)
											}
										} else if r[i] == "equipmentId" {
											if id, err := strconv.Atoi(r[i+1]); err == nil {
												result.EquipmentInfo.EquipmentID = int64(id)
											}
										}
									}
								}

								switch msg.DataType {
								case "float64":
									{
										if v, ok := msg.Value.(float64); ok {
											if result.DeviceType == "设备档案" {
												result.EquipmentInfo.Value = v
											} else if result.DeviceType == "环境档案" {
												result.EnvironmentAlarmInfo.Value = v
											}
										}
									}
								case "float32":
									{
										if v, ok := msg.Value.(float32); ok {
											if result.DeviceType == "设备档案" {
												result.EquipmentInfo.Value = float64(v)
											} else if result.DeviceType == "环境档案" {
												result.EnvironmentAlarmInfo.Value = float64(v)
											}
										}
									}
								case "uint32":
									{
										if v, ok := msg.Value.(uint32); ok {
											if result.DeviceType == "设备档案" {
												result.EquipmentInfo.Value = float64(v)
											} else if result.DeviceType == "环境档案" {
												result.EnvironmentAlarmInfo.Value = float64(v)
											}
										}
									}
								case "int32":
									{
										if v, ok := msg.Value.(int32); ok {
											if result.DeviceType == "设备档案" {
												result.EquipmentInfo.Value = float64(v)
											} else if result.DeviceType == "环境档案" {
												result.EnvironmentAlarmInfo.Value = float64(v)
											}
										}
									}
								}

								// 以-切割字符串
								// 获取 deviceType-xxx-equimentId-xxx-thresholdId-xxx-sensorName-xxx-value-xxx
								// if nodeModel.ID != 0 && nodeModel.Param != "" {
								// 	deviceTypeStart := strings.Index(nodeModel.Param, "deviceType-") + len("deviceType-")
								// 	deviceTypeEnd := strings.Index(nodeModel.Param[deviceTypeStart:], "-") + deviceTypeStart
								// 	// deviceType := str[deviceTypeStart:deviceTypeEnd]
								// 	if deviceTypeEnd > deviceTypeStart {
								// 		result.DeviceType = nodeModel.Param[deviceTypeStart:deviceTypeEnd]
								// 	}

								// 	environment := EnvironmentAlarmInfoDTO{}
								// 	environmentStart := strings.Index(nodeModel.Param, "environmentId-") + len("environmentId-")
								// 	environmentEnd := strings.Index(nodeModel.Param[environmentStart:], "-") + environmentStart
								// 	if environmentEnd > environmentStart {
								// 		environmentIdStr := nodeModel.Param[environmentStart:environmentEnd]
								// 		str, err := strconv.ParseInt(environmentIdStr, 10, 64)
								// 		if err == nil {
								// 			environment.EnvironmentID = str
								// 			switch msg.DataType {
								// 			case "float64":
								// 				{
								// 					environment.Value = msg.Value.(float64)
								// 				}
								// 			case "float32":
								// 				{
								// 					environment.Value = float64(msg.Value.(float32))
								// 				}
								// 			case "uint32":
								// 				{
								// 					environment.Value = float64(msg.Value.(uint32))
								// 				}
								// 			}
								// 		}
								// 	}

								// 	threhold := EquipmentInfoDTO{}
								// 	threholdStart := strings.Index(nodeModel.Param, "thresholdId-") + len("thresholdId-")
								// 	threholdEnd := strings.Index(nodeModel.Param[threholdStart:], "-") + threholdStart

								// 	if threholdEnd > threholdStart {
								// 		thresholdIdStr := nodeModel.Param[threholdStart:threholdEnd]
								// 		str, err := strconv.ParseInt(thresholdIdStr, 10, 64)
								// 		if err == nil {
								// 			threhold.ThresholdID = str
								// 		}
								// 		switch msg.DataType {
								// 		case "float64":
								// 			{
								// 				threhold.Value = msg.Value.(float64)
								// 			}
								// 		case "float32":
								// 			{
								// 				threhold.Value = float64(msg.Value.(float32))
								// 			}
								// 		case "uint32":
								// 			{
								// 				threhold.Value = float64(msg.Value.(uint32))
								// 			}
								// 		}

								// 		threhold.Value = msg.Value.(float64)
								// 		equipmentStart := strings.Index(nodeModel.Param, "equipment-") + len("equipment-")
								// 		equipmentEnd := strings.Index(nodeModel.Param[equipmentStart:], "-") + equipmentStart
								// 		if equipmentEnd > equipmentStart {
								// 			equipmentIdStr := nodeModel.Param[equipmentStart:equipmentEnd]
								// 			str, err := strconv.ParseInt(equipmentIdStr, 10, 64)
								// 			if err == nil {
								// 				threhold.EquipmentID = str
								// 			}
								// 		}
								// 	}
								// 	result.EnvironmentAlarmInfo = environment
								// 	result.EquipmentInfo = threhold

								// }

								fmt.Println("发送数据到后台==222=============", result)
								jsonStr, err := json.Marshal(result)
								if err != nil {
									continue
								}
								fmt.Println("发送数据到后台===============", string(jsonStr))
								// fmt.Println("发送数据到后台", string(jsonStr))

								stompConn.Send(global.Config.Sw.Topic, "application/json", jsonStr)
							}
						case <-ctx.Done():
							{
								return
							}
						}
					}
				}
			}
		}

		// // 订阅主题
		// sub, err := stompConn.Subscribe(global.Config.Sw.Topic, stomp.AckAuto)
		// if err != nil {
		// 	log.Fatalf("Failed to subscribe to topic: %v", err)
		// }
		// defer sub.Unsubscribe()

		// log.Println("Subscribed to /topic/example")
		// d := DeviceDTO{}
		// d.DeviceID = 1
		// d.DeviceType = "test"
		// d.EquipmentInfo.EquipmentID = 1
		// d.EquipmentInfo.SensorName = "test"
		// d.EquipmentInfo.ThresholdID = 1
		// d.EquipmentInfo.Value = 1
		// d.EnvironmentAlarmInfo.EnvironmentID = 1
		// jsonStr, _ := json.Marshal(d)
		// stompConn.Send(global.Config.Sw.Topic, "application/json", jsonStr)

	}()

}
