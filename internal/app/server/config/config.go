package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Server struct {
		Host   string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
		Port   int    `yaml:"port" env:"SERVER_PORT" env-default:"8081"`
		Secure bool   `yaml:"secure" env:"SERVER_SECURE" env-default:"true"`
	} `json:"server" yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
		Port     int    `yaml:"port" env:"PORT" env-default:"5432"`
		Name     string `yaml:"name" env:"NAME" env-default:"gophkeeper"`
		User     string `yaml:"username"  env-default:"postgres"`
		Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	} `yaml:"database"`
}

func MustLoad() *ServerConfig {
	var cfg ServerConfig
	if err := cleanenv.ReadConfig("config/server.yaml", &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func (c *ServerConfig) GetRepoURL() string {
	if c.Database.Host == "" {
		return ""
	}

	db := c.Database
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.User, db.Password, db.Host, db.Port, db.Name)
}

func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *ServerConfig) IsServerSecure() bool {
	return c.Server.Secure
}
