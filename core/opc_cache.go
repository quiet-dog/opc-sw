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

					global.DB.Where("id = ?", msg.ID).First(&node)
					jsonByte, err := json.Marshal(node)
					if err != nil {
						continue
					}
					id := fmt.Sprintf("%d-%s", node.ServiceId, node.NodeId)
					global.Redis.Set(global.Ctx, id, string(jsonByte), 0)
				}
			}
		}
	}()

}
