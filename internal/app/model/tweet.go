package model

import "time"

type Tweet struct {
	ID        uint      `gorm:"primary_key"`
	Author    string    `json:"author,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	URL       string    `json:"url,omitempty"`
}
