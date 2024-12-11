package main

import (
	"encoding/json"
	"log"
	"net/http"

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

func main() {
	// WebSocket 连接信息
	url := "ws://localhost:9020/ws"
	header := http.Header{}

	// 使用 gorilla/websocket 连接到服务器
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// 包装 WebSocket 连接为 io.ReadWriteCloser
	rwc := &WebSocketReadWriteCloser{Conn: conn}

	// 使用 STOMP 客户端连接
	stompConn, err := stomp.Connect(rwc)
	if err != nil {
		log.Fatalf("Failed to connect to STOMP: %v", err)
	}
	defer stompConn.Disconnect()

	log.Println("Connected to STOMP server")

	// 订阅主题
	sub, err := stompConn.Subscribe("/topic/info", stomp.AckAuto)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("Subscribed to /topic/example")
	d := DeviceDTO{}
	d.DeviceID = 1
	d.DeviceType = "test"
	d.EquipmentInfo.EquipmentID = 1
	d.EquipmentInfo.SensorName = "test"
	d.EquipmentInfo.ThresholdID = 1
	d.EquipmentInfo.Value = 1
	d.EnvironmentAlarmInfo.EnvironmentID = 1
	jsonStr, _ := json.Marshal(d)
	stompConn.Send("/app/deviceInfo", "application/json", jsonStr)

	// // 接收消息
	// for msg := range sub.C {
	// 	log.Printf("Received message: %s\n", string(msg.Body))
	// }
}
