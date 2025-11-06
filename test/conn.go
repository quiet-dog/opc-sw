package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

func main() {
	// 命令行参数
	endpoint := flag.String("ip", "opc.tcp://127.0.0.1:4840", "OPC UA 服务器地址")
	nodeID := flag.String("node", "ns=2;s=Demo.Static.Scalar.Float", "要读取的节点ID")
	flag.Parse()

	ctx := context.Background()

	// 创建 OPC UA 客户端
	client, err := opcua.NewClient(*endpoint,
		opcua.SecurityMode(ua.MessageSecurityModeNone),
		opcua.SecurityPolicy(ua.SecurityPolicyURINone),
		opcua.AutoReconnect(true),
	)
	if err != nil {
		panic(err)
	}

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close(ctx)

	fmt.Printf("✅ 已连接: %s\n", *endpoint)

	// 解析节点ID
	nid, err := ua.ParseNodeID(*nodeID)
	if err != nil {
		log.Fatalf("❌ 节点ID解析失败: %v", err)
	}

	// 读取节点值
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	val, err := client.Node(nid).Value(ctx)
	if err != nil {
		log.Fatalf("❌ 读取节点值失败: %v", err)
	}

	fmt.Printf("📦 节点 %s 的值: %v\n", *nodeID, val.Value())
}
