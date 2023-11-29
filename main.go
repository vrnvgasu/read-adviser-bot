package main

import (
	"flag"
	"log"
)

func main() {
	// токен получаем из параметров при запуске программы (чтобы не светить в коде)
	token := mustToken()

	// tgClient - через него идет общение с телегой
	// tgClient = tgClient.New(token)

	// fetcher - получает; processor - обрабатывает
	// fetcher отправляет в телегу сообщения, чтобы получать новые события
	// fetcher = fetcher.New()
	// processor отправляет в телегу новые сообщения после обработки
	// processor = processor.New()

	// consumer получает и обрабатывает сообщения.
	// consumer.Start(fetcher, processor)
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
