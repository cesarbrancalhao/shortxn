package domain

import "time"

type ClickEvent struct {
	URLId     string    `json:"url_id"`
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
}
