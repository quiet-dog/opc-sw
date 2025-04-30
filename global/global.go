package global

import (
	"context"
	"sw/config"
	"sw/opc"
	"sync"

	"github.com/lxzan/gws"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DEVICEDATA = 0
	YUZHI      = 1
	BAOJING    = 2
)

type RecHandler struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

var (
	DB         *gorm.DB
	Config     config.Config
	OpcGateway = opc.New()
	Ctx        = context.Background()
	Redis      *redis.Client
	Upgrader   *gws.Upgrader
	Session    = sync.Map{}
	RecChanel  = make(chan RecHandler, 5)
)
