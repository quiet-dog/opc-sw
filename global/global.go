package global

import (
	"sw/config"
	"sw/opc"

	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	Config     config.Config
	OpcGateway = opc.New()
)
