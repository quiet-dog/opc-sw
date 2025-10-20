package core

import (
	"fmt"
	"log"
	"sw/global"
	"sw/model/node"
	"sw/model/service"
	"sw/opc"
	"time"
)

func InitOpc() {
	fmt.Println(len(global.Config.Ingorenodes))

	var service []*service.ServiceModel
	global.DB.Find(&service)
	log.Println("初始化opc")
	for _, s := range service {
		log.Println("遍历服务", s.Opc)

		var nodes []*node.NodeModel
		global.DB.Where("service_id = ?", s.ID).Find(&nodes)
		var opcNodes []opc.NodeId
		for _, n := range nodes {
			isNotExit := false
			for _, not := range global.Config.Ingorenodes {
				if int64(n.ID) == not {
					isNotExit = true
				}
			}
			if !isNotExit {
				opcNodes = append(opcNodes, opc.NodeId{
					Node: n.NodeId,
					ID:   uint64(n.ID),
				})
			}
		}

		opcIP := s.Opc

		err := global.OpcGateway.AddClinet(fmt.Sprintf("%d", s.ID), opc.OpcClient{
			Endpoint: opcIP,
			Duration: time.Second * 60000,
			Nodes:    opcNodes,
			Username: s.Username,
			Password: s.Password,
		})
		if err != nil {
			fmt.Println("连接OPC服务器失败" + s.Opc)
			continue
		}

	}
	log.Println("初始化opc完成")

}
