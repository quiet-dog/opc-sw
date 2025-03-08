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
					global.DB.Where("id = ?", msg.ID).First(&node)
					jsonByte, err := json.Marshal(msg)
					if err != nil {
						fmt.Println("错误了")
						continue
					}
					msg.Param = node.Param
					id := fmt.Sprintf("%d", msg.ID)
					global.Redis.Set(global.Ctx, id, string(jsonByte), 0)
				}
			}
		}
	}()

}
