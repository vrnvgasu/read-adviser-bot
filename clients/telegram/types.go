package telegram

// https://core.telegram.org/bots/api#making-requests
type UpdatesResponse struct {
	Ok     string   `json:"ok"`
	Result []Update `json:"result"`
}

// структура ответа из документации https://core.telegram.org/bots/api#update
type Update struct {
	ID      int    `json:"update_id"` // json будет парсить update_id на ID
	Message string `json:"message"`
}
