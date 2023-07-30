package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host        string `required:"true" envconfig:"HOST"`
	Port        int    `required:"true" envconfig:"PORT"`
	Username    string `required:"true" envconfig:"USERNAME"`
	Password    string `required:"true" envconfig:"PASSWORD"`
	Database    string `required:"true" envconfig:"DATABASE"`
	QueryString string `required:"true" envconfig:"QUERYSTRING"`
}

func (c *DatabaseConfig) RWDataSourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.QueryString,
	)
}

func NewDB(dbCfg DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dbCfg.RWDataSourceName()), &gorm.Config{})
}

func loadGatewayConfig() (GatewayConfig, error) {
	var cfg GatewayConfig

	// load from .env if exists
	if _, err := os.Stat(".env"); err == nil {
		if err := gotenv.Load(); err != nil {
			return cfg, err
		}
	}

	// parse environment variable to config struct using "service" namespace
	// to prevent conflict with another modules
	err := envconfig.Process("service", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
