package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/gsk148/gophkeeper/internal/app/client"
)

func main() {
	c, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if err = c.Register("test", "super"); err != nil {
		log.Error(err)
	}

	if err = c.Login("test", "super"); err != nil {
		log.Error(err)
	}

	id, err := c.StorePassword("test", "test", "test", "note")
	if err != nil {
		log.Error(err)
	}
	log.Info(id)

	bin, err := c.GetPasswordByID(id)
	if err != nil {
		log.Error(err)
	}
	log.Info(bin)

	bins, err := c.GetAllPasswords()
	if err != nil {
		log.Error(err)
	}
	log.Info(bins)

	if err = c.DeletePassword(id); err != nil {
		log.Error(err)
	}

	if err = c.Logout(); err != nil {
		log.Error(err)
	}
}
