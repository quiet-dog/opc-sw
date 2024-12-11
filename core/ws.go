package core

import (
	"encoding/json"
	"sw/global"
	"sw/opc"
	"sync"
	"time"

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

func (c *Handler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	notifyChan := global.OpcGateway.SubscribeOpc()
	global.Session.Store(socket, notifyChan)
	// 获取redis缓存的所有数据
	keys, _ := global.Redis.Keys(global.Ctx, "*").Result()
	for _, key := range keys {
		var notify opc.Notify
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

	go func() {
		for {
			select {
			case msg, ok := <-notifyChan:
				{
					if !ok {
						return
					}
					b, err := json.Marshal(msg)
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
		if notify, ok := v.(chan opc.Notify); ok {
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
