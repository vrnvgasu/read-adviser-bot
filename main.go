package main

import (
	"flag"
	"log"
	client "read-adviser-bot/clients/telegram"
	event_consumer "read-adviser-bot/consumer/event-consumer"
	"read-adviser-bot/events/telegram"
	"read-adviser-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	// tgClient - через него идет общение с телегой
	tgClient := client.New(tgBotHost, mustToken())
	eventsProcessor := telegram.New(tgClient, files.New(storagePath))

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
