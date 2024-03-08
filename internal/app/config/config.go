package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DatabaseConfig `yaml:"database"`
	ServerConfig   `yaml:"http_server""`
}

type ServerConfig struct {
	Addr string `yaml:"addr" env:"ADDR" env-default:"localhost:8081"`
}

type DatabaseConfig struct {
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"HOST" env-default:"127.0.0.1"`
	Name     string `yaml:"name" env:"NAME" env-default:"gophkeeper"`
	Username string `yaml:"username"  env-default:"postgres"`
	Password string `yaml:"password" env:"postgres"`
	Sslmode  string `yaml:"sslmode" env-default:"disable"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig("config/config.yaml", &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
