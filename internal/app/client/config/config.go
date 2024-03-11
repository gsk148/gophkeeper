package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"

	"github.com/gsk148/gophkeeper/internal/pkg/cert"
)

type ClientConfig struct {
	API struct {
		Host   string `yaml:"host" env:"API_HOST" env-default:"localhost"`
		Port   int    `yaml:"port" env:"API_PORT" env-default:"8081"`
		Route  string `yaml:"route" env:"API_ROUTE" env-default:"/api/v1"`
		Secure bool   `yaml:"secure" env:"CLIENT_SECURE" env-default:"false"`
	} `yaml:"api"`
	Cert struct {
		CA   string `yaml:"ca" env:"CA_PATH" env-default:"cert/CertAuth.crt"`
		Cert string `yaml:"cert" env:"CLIENT_CERT_PATH" env-default:"cert/client.crt"`
		Key  string `yaml:"key" env:"CLIENT_KEY_PATH" env-default:"cert/client.key"`
	} `yaml:"cert"`
}

func MustLoad() *ClientConfig {
	var cfg ClientConfig
	if err := cleanenv.ReadConfig("/home/anton/dev/gophkeeper/config/client.yaml", &cfg); err != nil {
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

func (c *ClientConfig) GetCACertPool() (*x509.CertPool, error) {
	return cert.GetCertificatePool(c.Cert.Cert)
}

func (c *ClientConfig) GetCertificate() (tls.Certificate, error) {
	return cert.GetClientCertificate(c.Cert.Cert, c.Cert.Key)
}
