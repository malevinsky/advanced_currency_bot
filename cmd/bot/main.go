package main

import (
	"context"
	"errors"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/clients/tg"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/config"
	"gitlab.ozon.dev/amalevinskaya/teodora-malevinskaia/internal/model/messages"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	config, err := config.New()
	if err != nil {
		err := errors.New("config init failed:")
		if err != nil {
			return
		}
	}

	tgClient, err := tg.New(config)
	if err != nil {
		errors.New("tg client init failed:")
	}

	msgModel := messages.New(tgClient)

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "allDoneWG", &sync.WaitGroup{})

	tgClient.ListenUpdates(msgModel)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		if err := recover(); err != nil {
			log.Println("recovered from panic", err)
		}
		<-exit
		cancel()
	}()
}
