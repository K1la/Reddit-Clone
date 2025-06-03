package models

import "time"

type (
	Comment struct {
		Created time.Time `json:"created"`
		Author  *Session  `json:"author"`
		Body    string    `json:"body"`
		ID      string    `json:"id"`
	}
)
