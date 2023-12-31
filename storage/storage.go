package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"read-adviser-bot/lib/e"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error) // кому показать
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error) // существует ли страница
}

// страница, на которую ведет ссылка, которую скинули боту
type Page struct {
	URL      string
	UserName string // какому пользователю возвращать
}

var ErrNoSavedPages = errors.New("no saved pages")

func (p Page) Hash() (string, error) {
	h := sha1.New()
	// h - hash.Hash - реализует интерфейс Writer
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	// h.Sum(nil) - сложили хеш урла и имени пользователя (посчитали их выше)
	// "%x" преобразовать байты в строку (h.Sum(nil) - вернет байты)
	result := fmt.Sprintf("%x", h.Sum(nil))

	return result, nil
}
