package dto

type WebHookDTO struct {
	URL        string      `json:"url"`
	User_id    string      `json:"user_id"`
	Payload    interface{} `json:"payload"`
	RetryCount int         `json:"retry_count"`
}
