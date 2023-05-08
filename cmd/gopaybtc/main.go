package main

import (
	"github.com/ngenohkevin/go-pay-btc/pkg/gopaybtc"
	"log"
)

func main() {

	config, err := gopaybtc.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	server, err := gopaybtc.NewServer(config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("unable to start the server", err)
	}

}
