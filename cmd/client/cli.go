package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/gsk148/gophkeeper/internal/app/cli"
)

func main() {
	client, err := cli.NewCLI()
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Start(); err != nil {
		log.Fatal(err)
	}
}
