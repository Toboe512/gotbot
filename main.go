package main

import (
	"flag"
	tgClient "github.com/toboe512/gotbot/clients/telegram"
	"github.com/toboe512/gotbot/consumer/event-consumer"
	"github.com/toboe512/gotbot/events/telegram"
	"github.com/toboe512/gotbot/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mastToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(
		eventsProcessor,
		eventsProcessor,
		batchSize,
	)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped")
	}
}

func mastToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
