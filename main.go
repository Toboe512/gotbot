package main

import (
	"context"
	"flag"
	tgClient "github.com/toboe512/gotbot/clients/telegram"
	"github.com/toboe512/gotbot/consumer/event-consumer"
	"github.com/toboe512/gotbot/events/telegram"
	"github.com/toboe512/gotbot/storage"
	"github.com/toboe512/gotbot/storage/sqlite"
	"github.com/toboe512/gotbot/utils"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	storageFilePath   = "files_storage"
	storageSqlitePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	ctx := context.TODO()
	tkn := mastToken()
	str := getSqliteStorage(ctx, storageSqlitePath)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, tkn),
		str,
	)

	log.Print("service started")

	consumer := event_consumer.New(
		eventsProcessor,
		eventsProcessor,
		batchSize,
	)

	if err := consumer.Start(ctx); err != nil {
		log.Fatal("service is stopped")
	}
}

func mastToken() string {
	token := flag.String(
		"tg-bot-token",
		utils.EmptyStr,
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == utils.EmptyStr {
		log.Fatal("token is not specified")
	}

	return *token
}

//func getFileStorage(path string) storage.Storage {
//	return files.New(storageFilePath)
//}

func getSqliteStorage(ctx context.Context, path string) storage.Storage {
	s, err := sqlite.New(path)

	if err != nil {
		log.Fatal("can't connect to sqlite: %w", err)
	}

	if err := s.Init(ctx); err != nil {
		log.Fatal("can't init to sqlite: %w", err)
	}

	return s
}
