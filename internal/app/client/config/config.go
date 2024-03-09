package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type ClientConfig struct {
	API struct {
		Host   string `yaml:"host" env:"API_HOST" env-default:"localhost"`
		Port   int    `yaml:"port" env:"API_PORT" env-default:"8081"`
		Route  string `yaml:"route" env:"API_ROUTE" env-default:"/api/v1"`
		Secure bool   `yaml:"secure" env:"CLIENT_SECURE" env-default:"false"`
	} `yaml:"api"`
}

func MustLoad() *ClientConfig {
	var cfg ClientConfig
	if err := cleanenv.ReadConfig("config/client.yaml", &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func (c *ClientConfig) GetAPIAddress() string {
	protocol := "http"
	if c.API.Secure {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d%s", protocol, c.API.Host, c.API.Port, c.API.Route)
}
