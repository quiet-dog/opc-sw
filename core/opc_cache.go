package core

import (
	"encoding/json"
	"fmt"
	"sw/global"
	"sw/model/node"
)

func InitOpcCache() {
	notify := global.OpcGateway.SubscribeOpc()

	go func() {
		for {
			select {
			case msg, ok := <-notify:
				{
					if !ok {
						fmt.Println("opc cache notify channel closed")
						return
					}

					var node node.NodeModel
					fmt.Println("存入缓存")

					global.DB.Where("node_id = ?", msg.NodeId).First(&node)
					msg.Params = node.Param
					jsonByte, err := json.Marshal(msg)
					if err != nil {
						continue
					}
					global.Redis.Set(global.Ctx, msg.NodeId, string(jsonByte), 0)
				}
			}
		}
	}()

}
