package telegram

import "read-adviser-bot/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
	// storage - сюда будем сохранять ссылки
}

func New(client *telegram.Client) {

}
