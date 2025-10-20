package opc

import (
	"context"
	"errors"
	"sync"
	"time"
)

type OpcGateway struct {
	opcs   sync.Map
	notify chan Data
	sub    sync.Map
}
type Config struct {
	Endpoint string
	Duration time.Duration
	ctx      context.Context
}

func New() *OpcGateway {
	o := &OpcGateway{}
	o.notify = make(chan Data)

	go func() {
		for {
			select {
			case msg, ok := <-o.notify:
				if ok {
					o.sub.Range(func(key, value interface{}) bool {
						ch := key.(chan Data)
						select {
						case ch <- msg:
							// 成功发送
						default:
							// 没人接收，跳过
						}
						return true
					})
					continue
				}
				return
			}
		}
	}()
	return o
}

func (o *OpcGateway) AddClinet(clientId string, config OpcClient) error {
	c := &OpcClient{
		Endpoint: config.Endpoint,
		Duration: config.Duration,
		gateway:  o.notify,
		Nodes:    config.Nodes,
		Username: config.Username,
		Password: config.Password,
	}

	go c.connect()
	o.opcs.Store(clientId, c)
	return nil
}

func (o *OpcGateway) AddNode(clientId string, nodeId NodeId) error {
	c, ok := o.opcs.Load(clientId)
	if !ok {
		return errors.New("client not found")
	}
	client := c.(*OpcClient)
	client.AddNodeID(nodeId)
	return nil
}

// 订阅
func (o *OpcGateway) SubscribeOpc() <-chan Data {
	ch := make(chan Data)
	o.sub.Store(ch, nil)
	return ch
}

// 取消订阅
func (o *OpcGateway) UnSubscribeOpc(ch <-chan Data) {
	o.sub.Delete(ch)
}
