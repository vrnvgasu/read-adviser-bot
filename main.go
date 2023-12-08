package main

import (
	"context"
	"flag"
	"log"
	client "read-adviser-bot/clients/telegram"
	event_consumer "read-adviser-bot/consumer/event-consumer"
	"read-adviser-bot/events/telegram"
	"read-adviser-bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	//storage := files.New(storagePath) // хранение в файлах
	storage, err := sqlite.New(sqliteStoragePath) // хранение в БД
	if err != nil {
		log.Fatal("can't connect to storage", err)
	}

	// `TODO` - указываем, что еще не определились, какой конкретно контекст будем использовать
	if err := storage.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage", err)
	}

	// tgClient - через него идет общение с телегой
	tgClient := client.New(tgBotHost, mustToken())
	eventsProcessor := telegram.New(tgClient, storage)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

// must - указание на то, что с этой функцией надо быть осторожным
// бросит панику при ошибке
func mustToken() string {
	// запускаем программу с: bot -tg-bot-token 'my token'
	// сюда попадет только ссылка на значение
	token := flag.String("tg-bot-token", "", "token to access to telegram bot")
	flag.Parse() // только в этот момент в token запишется ссылка на значение

	if *token == "" {
		log.Fatal("token is not provided")
	}

	return *token
}
