package core

import (
	"context"
	"encoding/json"
	"fmt"
	"sw/global"
	"sw/model/node"
	"time"
)

func InitOpcCache() {
	notify := global.OpcGateway.SubscribeOpc()

	go func() {
		for {
			select {
			case <-context.Background().Done():
				return
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
						fmt.Println("错误了", err)
						continue
					}
					msg.Param = node.Param
					id := fmt.Sprintf("%d", msg.ID)
					global.Redis.Set(global.Ctx, id, string(jsonByte), 5*time.Second)
				}
			}
		}
	}()

	go func() {
		for {
			ctx := context.Background()
			select {
			case <-ctx.Done():
				fmt.Println("context done, exiting opc cache goroutine")
				return
			case msg, ok := <-global.MoniChannel:
				{
					if !ok {
						fmt.Println("monitor channel closed")
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
