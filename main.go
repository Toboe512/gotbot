package main

import (
	"flag"
	tgClient "goTBot/clients/telegram"
	"goTBot/consumer/event-consumer"
	"goTBot/events/telegram"
	"goTBot/storage/files"
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
