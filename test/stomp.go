package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/gorilla/websocket"
)

// WebSocketWrapper 包装 websocket.Conn 以实现 io.ReadWriteCloser
type WebSocketWrapper struct {
	conn *websocket.Conn
}

func (w *WebSocketWrapper) Read(p []byte) (n int, err error) {
	// 使用 websocket.Conn 的 ReadMessage 来实现 io.Reader
	_, p, err = w.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *WebSocketWrapper) Write(p []byte) (n int, err error) {
	// 使用 websocket.Conn 的 WriteMessage 来实现 io.Writer
	err = w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *WebSocketWrapper) Close() error {
	return w.conn.Close()
}

func main() {
	// 连接 WebSocket 服务器
	wsConn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9020/ws", nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server: ", err)
	}
	defer wsConn.Close()
	fmt.Println("=323232")

	// 将 WebSocket 连接包装为实现了 io.ReadWriteCloser 接口的类型
	wrappedConn := &WebSocketWrapper{conn: wsConn}

	// 使用封装后的连接连接到 STOMP 协议
	conn, err := stomp.Connect(wrappedConn)
	fmt.Println("=132123323232")
	if err != nil {
		log.Fatal("Failed to connect with STOMP protocol: ", err)
	}
	defer conn.Disconnect()
	fmt.Println("=12313")
	// 订阅消息
	subscription, err := conn.Subscribe("/topic/info", stomp.AckAuto)
	if err != nil {
		log.Fatal("Failed to subscribe: ", err)
	}
	defer subscription.Unsubscribe()

	// 启动一个goroutine处理接收到的消息
	go func() {
		for {
			// 接收消息
			message, err := subscription.Read()
			if err != nil {
				log.Fatal("Failed to read message: ", err)
			}
			fmt.Printf("Received message: %s\n", string(message.Body))
		}
	}()

	fmt.Println("=到爱上对方")
	// 发送消息到Spring Boot服务器
	err = conn.Send("/app/deviceInfo", "text/plain", []byte("Hello from Go"))
	if err != nil {
		log.Fatal("Failed to send message: ", err)
	}
	fmt.Println("Message sent to Spring Boot server")

	// 防止主程序退出
	time.Sleep(10 * time.Second)
}
