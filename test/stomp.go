package main

import (
	"fmt"
	"log"

	"github.com/go-stomp/stomp/v3"
	"github.com/gorilla/websocket"
)

// WebSocketWrapper 包装 websocket.Conn 以实现 io.ReadWriteCloser
type WebSocketWrapper struct {
	conn *websocket.Conn
}

func (w *WebSocketWrapper) Read(p []byte) (n int, err error) {
	_, p, err = w.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *WebSocketWrapper) Write(p []byte) (n int, err error) {
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
	// Endereço do servidor STOMP
	serverAddress := "localhost:9020"

	// Conectar ao servidor STOMP
	conn, err := stomp.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()

	// Destino para enviar mensagens
	destination := "/topic/greetings"

	// Corpo da mensagem
	messageBody := "Hello, STOMP!"

	// Subscrever para receber mensagens
	sub, err := conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		println(err)
	}
	defer sub.Unsubscribe()

	// Enviar uma mensagem
	err = conn.Send(destination, "text/plain", []byte(messageBody), nil)
	if err != nil {
		println(err)
	}

	// Aguardar por mensagens
	for {
		msg := <-sub.C
		fmt.Printf("Mensagem recebida: %s", msg.Body)
	}
}
