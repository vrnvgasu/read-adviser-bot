package telegram

// https://core.telegram.org/bots/api#making-requests
type UpdatesResponse struct {
	Ok     string   `json:"ok"`
	Result []Update `json:"result"`
}

// структура ответа из документации https://core.telegram.org/bots/api#update
type Update struct {
	ID      int              `json:"update_id"` // json будет парсить update_id на ID
	Message *IncomingMessage `json:"message"`   // типа, если ссылка на структуру *, то поле опциональное, может не прилететь
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
