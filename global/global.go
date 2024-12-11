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
	DB         *gorm.DB
	Config     config.Config
	OpcGateway = opc.New()
	Ctx        = context.Background()
	Redis      *redis.Client
	Upgrader   *gws.Upgrader
	Session    = sync.Map{}
)
