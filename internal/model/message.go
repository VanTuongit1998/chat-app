package model

import "time"

type Message struct {
	Sender  string    `json:"sender"`
	Message string    `json:"message"`
	To      string    `json:"to,omitempty"`
	SentAt  time.Time `json:"sentAt"`
}
