package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// WebSocket 地址
	url := "ws://127.0.0.1:9180/ws"

	// 连接 WebSocket
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer c.Close()
	log.Println("已连接到 WebSocket 服务器")

	// 捕获中断信号，优雅退出
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// 接收消息的协程
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("读取错误:", err)
				return
			}
			log.Printf("收到消息: %s", message)
		}
	}()

	// 发送心跳或示例消息
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			// 发送心跳消息
			err := c.WriteMessage(websocket.TextMessage, []byte("ping "+t.Format(time.RFC3339)))
			if err != nil {
				log.Println("发送错误:", err)
				return
			}
		case <-interrupt:
			log.Println("接收到中断信号，关闭连接...")
			// 发送关闭消息
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("关闭错误:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
