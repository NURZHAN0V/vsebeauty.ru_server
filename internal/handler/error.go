package handler

// ErrorResponse — стандартный формат ошибки API
type ErrorResponse struct {
	Error   string `json:"error"`             // Сообщение об ошибке
	Details string `json:"details,omitempty"` // Дополнительные детали (необязательно)
}
