package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage"
	"time"
)

type Storage struct {
	basePath string // в какой папке храним
}

const defaultPermission = 0774 // права на чтение и запись всем пользователям

var ErrNoSavedPages = errors.New("no saved pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	// перехватываем вконце ошибку и оборочиваем ее.
	// чтобы сработало надо сделать возвращаем тип именной во всем методе
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	// filepath.Join создает правильный путь на линуксу и винде
	filePath := filepath.Join(s.basePath, page.UserName)

	// создаем путь
	if err := os.MkdirAll(filePath, defaultPermission); err != nil {
		return err
	}

	fileName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	// закрываем ресурс в конце (без обработки ошибки
	defer func() { _ = file.Close() }()

	// преобразуем страницу в формат gob и записываем в file
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	// список файлов в каталоге
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// ошибка, если каталог пустой
	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	// rand - псевдослучайные числа. Нужно указать seed (тут время), чтобы при запуске генерились новые значения
	rand.NewSource(time.Now().UnixNano())
	n := rand.Intn(len(files)) //случайное число с верхней границей

	file := files[n]
	// open decode - открыть файл (в формате gob) и декодировать его
	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}

	return nil
}

// сохранял ли пользователь эту страницу ранее
func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fileName, err := fileName(page)
	if err != nil {
		return false, e.Wrap("can't find file", err)
	}

	path := filepath.Join(s.basePath, fileName)

	// проверяет существование (возвращает несколько вариантов ошибок)
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, err
	case err != nil:
		msg := fmt.Sprintf("can't find file %s", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = file.Close() }()

	// сюда декодируем страницу из файла gob
	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
