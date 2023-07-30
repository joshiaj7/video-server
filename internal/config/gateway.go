package config

import (
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"video-server/module/config"
	"video-server/module/entity"
)

type GatewayConfig struct {
	Environment    string         `envconfig:"ENVIRONMENT" default:"dev"`
	DatabaseConfig DatabaseConfig `envconfig:"DB"`

	Database *gorm.DB           `ignored:"true"`
	Router   *httprouter.Router `ignored:"true"`
}

func NewGatewayServer() (GatewayConfig, error) {
	cfg, err := loadGatewayConfig()
	if err != nil {
		return cfg, err
	}

	// init DB
	cfg.Database, err = NewDB(cfg.DatabaseConfig)
	if err != nil {
		return cfg, err
	}

	// migrate DB
	migrateDB(cfg.Database)

	// init router
	cfg.Router = httprouter.New()

	// register module
	moduleRepo := config.RegisterRepository(cfg.Database)
	moduleUsecase := config.RegisterUsecase(moduleRepo)
	config.RegisterHandler(cfg.Router, moduleUsecase)

	return cfg, nil
}

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&entity.File{})
}
