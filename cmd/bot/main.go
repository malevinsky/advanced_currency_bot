package main

import (
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/clients/tg"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/config"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/model/messages"
	"log"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(config)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	msgModel := messages.New(tgClient)

	tgClient.ListenUpdates(msgModel)
}
