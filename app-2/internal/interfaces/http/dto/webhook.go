package dto

type WebhookPayload struct {
	URL     string      `json:"url"`
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}
