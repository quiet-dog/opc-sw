package core

import (
	"fmt"
	"log"
	"sw/global"

	"github.com/go-stomp/stomp/v3"
)

func InitSw() {
	conn, err := stomp.Dial("tcp", fmt.Sprintf("%s:%s", global.Config.Sw.Host, global.Config.Sw.Port))

	if err != nil {
		log.Fatalf("无法连接到 STOMP 服务器: %v", err)
	}
	defer conn.Disconnect()

	// 订阅一个队列
	sub, err := conn.Subscribe(global.Config.Sw.Topic, stomp.AckAuto)
	if err != nil {
		log.Fatalf("无法订阅: %v", err)
	}
	defer sub.Unsubscribe()

	// 发送消息到 Spring Boot WebSocket
	err = conn.Send("/app/hello", "text/plain", []byte("Hello from Go client"))
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	// 接收并打印来自服务器的消息
	for {
		message := <-sub.C
		fmt.Printf("接收到消息: %s\n", message.Body)
	}
}
