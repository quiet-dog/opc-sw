package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"sw/global"
	"sw/model/node"
	"sw/opc"
	"sync"
	"time"

	"math/rand"

	"github.com/lxzan/gws"
)

const (
	PingInterval = 10000000 * time.Second
	PingWait     = 10000000 * time.Second
)

type Handler struct{}

type Session struct {
	Seection sync.Map
}

func getCustomRandom() int {
	rand.Seed(time.Now().UnixNano())

	// 决定从哪个区间选：true 表示选第一段，false 表示选第二段
	if rand.Intn(2) == 0 {
		// 第一段：[6, 99]
		return rand.Intn(94) + 6 // 99 - 6 + 1 = 94
	} else {
		// 第二段：[151, 404]
		return rand.Intn(254) + 151 // 404 - 151 + 1 = 254
	}
}

func (c *Handler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	notifyChan := global.OpcGateway.SubscribeOpc()
	global.Session.Store(socket, notifyChan)
	// 获取redis缓存的所有数据
	keys, _ := global.Redis.Keys(global.Ctx, "*").Result()
	for _, key := range keys {
		var notify opc.Data
		err := global.Redis.Get(global.Ctx, key).Scan(notify)
		if err != nil {
			continue
		}
		json, err := json.Marshal(notify)
		if err != nil {
			continue
		}
		socket.WriteMessage(gws.OpcodeText, json)
	}

	// 定时器
	// timer := time.NewTicker(5 * time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-timer.C:
	// 			{
	// 				notify := opc.Notify{}
	// 				notify.NodeId = "1"
	// 				notify.Value = "1"
	// 				jsonByte, _ := json.Marshal(notify)
	// 				socket.WriteMessage(gws.OpcodeText, jsonByte)
	// 			}
	// 		}
	// 	}
	// }()
	testchanel := make(chan opc.Data)

	go func() {
		{
			fmt.Println("=================")
			var nodeModel []node.NodeModel
			for _, v := range nodeModel {
				notify := opc.Data{}
				notify.ID = uint64(v.ID)
				notify.DataType = v.NodeId
				notify.SourceTime = time.Now()
				if strings.Contains(v.Key, "状态") || strings.Contains(v.Key, "报警") || strings.Contains(v.Key, "开关") {
					notify.Value = "1"
					rand.Seed(time.Now().UnixNano()) // 设置随机种子
					if rand.Intn(2) == 0 {
						notify.Value = true
					} else {
						notify.Value = false
					}
					jsonByte, _ := json.Marshal(notify)
					global.Redis.Set(global.Ctx, fmt.Sprintf("%d", v.ID), jsonByte, 0)
				} else {
					// 0-100的随机2位浮点数
					rand.Seed(time.Now().UnixNano()) // 设置随机种子
					// 生成 0 到 100 之间的随机 float64（保留两位小数）
					num := float64(rand.Intn(10000)) / 100.0
					notify.Value = num
					jsonByte, _ := json.Marshal(notify)
					global.Redis.Set(global.Ctx, fmt.Sprintf("%d", v.ID), jsonByte, 0)
				}
			}
		}
		for {
			fmt.Println("=================")
			randId := getCustomRandom()
			var v node.NodeModel
			global.DB.Where("id = ?", randId).Find(&v)

			// for _, v := range nodeModel {
			notify := opc.Data{}
			notify.ID = uint64(v.ID)
			notify.DataType = v.NodeId
			notify.SourceTime = time.Now()
			if strings.Contains(v.Key, "状态") || strings.Contains(v.Key, "报警") || strings.Contains(v.Key, "开关") {
				notify.Value = "1"
				// 随机生成true或false
				// 生成随机布尔值

				rand.Seed(time.Now().UnixNano()) // 设置随机种子
				notify.Value = rand.Intn(2) == 1 // 生成 0 或 1，判断是否为 1
				jsonByte, _ := json.Marshal(notify)
				global.Redis.Set(global.Ctx, fmt.Sprintf("%d", v.ID), jsonByte, 0)
			} else {
				// 0-100的随机2位浮点数
				rand.Seed(time.Now().UnixNano()) // 设置随机种子
				// 生成 0 到 100 之间的随机 float64（保留两位小数）
				num := float64(rand.Intn(10000)) / 100.0
				notify.Value = num
				jsonByte, _ := json.Marshal(notify)
				global.Redis.Set(global.Ctx, fmt.Sprintf("%d", v.ID), jsonByte, 0)
			}
			testchanel <- notify
			time.Sleep(1 * time.Second)
			// }
		}
	}()

	go func() {
		for {
			select {
			case msg, ok := <-notifyChan:
				{

					if !ok {
						return
					}
					result, err := getResult(msg)
					if err != nil {
						continue
					}
					b, err := json.Marshal(result)
					if err != nil {
						continue
					}
					socket.WriteMessage(gws.OpcodeText, b)
				}
			case msg, ok := <-testchanel:
				{
					if !ok {
						return
					}
					result, err := getResult(msg)
					if err != nil {
						continue
					}
					b, err := json.Marshal(result)
					if err != nil {
						continue
					}
					socket.WriteMessage(gws.OpcodeText, b)
				}
			}
		}
	}()
}

func (c *Handler) OnClose(socket *gws.Conn, err error) {
	if v, ok := global.Session.Load(socket); ok {
		// global.OpcGateway.UnsubscribeOpc(v)
		if notify, ok := v.(chan opc.Data); ok {
			close(notify)
		}
	}
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}

func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	socket.WriteMessage(message.Opcode, message.Bytes())
}

func InitWs() {
	upgrader := gws.NewUpgrader(&Handler{}, &gws.ServerOption{
		ParallelEnabled:   true,                                 // Parallel message processing
		Recovery:          gws.Recovery,                         // Exception recovery
		PermessageDeflate: gws.PermessageDeflate{Enabled: true}, // Enable compression
	})
	global.Upgrader = upgrader
}

func getResult(msg opc.Data) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	var 节点 node.NodeModel
	对应节点列表 := []*node.NodeModel{}

	err := global.DB.Where("id = ?", msg.ID).First(&节点).Error
	if err != nil {
		return nil, err
	}
	global.DB.Where("device_type = ?", 节点.DeviceType).Find(&对应节点列表)

	if 节点.DeviceType == "关键设备" {
		设备数据 := map[string]interface{}{}
		for _, 对应节点 := range 对应节点列表 {
			// 判断是否有对应的key
			if 设备数据[对应节点.DeviceName] == nil {
				父节点 := map[string]interface{}{}
				父节点[对应节点.Key] = 对应节点.Value
				设备数据[对应节点.DeviceName] = 父节点
			} else {
				设备数据[对应节点.DeviceName].(map[string]interface{})[对应节点.Key] = 对应节点.Value
			}
		}
		result["设备数据"] = 设备数据
	}

	if 节点.DeviceType == "EMS" {
		设备数据 := []map[string]interface{}{}
		for _, 对应节点 := range 对应节点列表 {
			var isExit bool
			for i, v := range 设备数据 {
				if v["区域"] == 对应节点.EmsAare {
					isExit = true
					设备数据[i][对应节点.Key] = 对应节点.Value
					break
				}
			}
			if !isExit {
				父节点 := map[string]interface{}{}
				父节点["区域"] = 对应节点.BmsArea
				父节点[对应节点.Key] = 对应节点.Value
				设备数据 = append(设备数据, 父节点)
			}
		}
		result["systemName"] = "EMS"
		result["Data"] = 设备数据
	}

	if 节点.DeviceType == "BMS" {
		设备数据 := map[string]interface{}{}
		for _, 对应节点 := range 对应节点列表 {
			if 设备数据[对应节点.BmsDeviceName] == nil {
				设备列表 := []map[string]interface{}{}
				设备信息 := map[string]interface{}{
					"区域":   对应节点.BmsArea,
					"设备标签": 对应节点.BmsLabel,
				}
				设备信息[对应节点.Key] = 对应节点.Value
				设备列表 = append(设备列表, 设备信息)
				设备数据[对应节点.DeviceName] = 设备列表
			} else {
				设备列表 := 设备数据[对应节点.BmsDeviceName].([]map[string]interface{})
				for i, 设备 := range 设备列表 {
					if 设备["区域"] == 对应节点.BmsArea {
						设备列表[i][对应节点.Key] = 对应节点.Value
						break
					}
				}
			}
		}
		for k, v := range 设备数据 {
			result[k] = v
		}
		result["systemName"] = "BMS"
	}

	return result, nil
}
