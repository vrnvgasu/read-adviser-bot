package telegram

import (
	"errors"
	"log"
	"net/url"
	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text) // удалим лишние пробелы

	log.Printf("get new command '%s' from '%s'", text, username)

	// add page: http:://... - просто сохранить ссылку
	// rnd page: /rnd - получить рандомную страницу
	// help: /help - объяснить, как работает бот
	// start: /start - автоматическая команда при начале общения с ботом

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	exists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if exists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can't do command: send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	// из минусов: ya.ru не будет считаться ссылкой
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
