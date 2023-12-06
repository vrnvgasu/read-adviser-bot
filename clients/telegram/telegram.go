package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"read-adviser-bot/lib/e"
	"strconv"
)

type Client struct {
	host     string      // хост api телеги
	basePath string      // путь для запросов host/bot<token>
	client   http.Client // чтобы не создавать клиента для каждого запроса
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// get не пишем
func newBasePath(token string) string {
	return "bot" + token
}

// getUpdates не пишем в go
// параметры запроса offset и limit взяли из доки https://core.telegram.org/bots/api#getupdates
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	// формируем запрос
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	// https://core.telegram.org/bots/api#getupdates - дока для метода телеги getUpdates
	// отправялем запрос на "getUpdates" и получаем ответ data
	data, err := c.doRequest(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	// парсим ответ из data
	var res UpdatesResponse
	// json.Unmarshal парсит data в структуру UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// в доке метод телеги для отправки сообщений "sendMessage"
// https://core.telegram.org/bots/api#sendmessage
// параметры chat_id и text из структуры запроса этой доки
func (c *Client) SendMessage(chatID int, text string) error {
	// формируем запрос
	query := url.Values{}
	query.Add("chatID", strconv.Itoa(chatID))
	query.Add("text", text)

	// отправляем запрос на "sendMessage"
	// тело ответа тут не интересно
	_, err := c.doRequest(sendMessageMethod, query)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

// возвращаемые параметры именнованные, чтобы сработал defer
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	// вызовется вконце и обернет err в e.WrapIfErr
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	// body - nill, тк все уже вставили в параметры query
	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// передаем параметры запроса в request
	request.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	// закрыли тело ответа в defer - закроет в конце метода
	defer func() { _ = resp.Body.Close() }()

	// читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
