package global

import (
	"sw/config"

	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Config config.Config
)
