package core

import (
	"log"
	"sw/global"
	"sw/model/node"
	"sw/model/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitOrm() {
	log.Println("初始化orm======================")
	db, err := gorm.Open(sqlite.Open(global.Config.Sqlite.Path), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&service.ServiceModel{},
		&node.NodeModel{},
	)
	global.DB = db
	log.Println("完成orm======================")

}
