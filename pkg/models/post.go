package models

import "time"

type (
	Post struct {
		Score      int        `json:"score"`
		Views      int        `json:"views"`
		Type       string     `json:"type"`
		Title      string     `json:"title"`
		Author     *Session   `json:"author"`
		Category   string     `json:"category"`
		Text       string     `json:"text,omitempty"`
		URL        string     `json:"url,omitempty"`
		Votes      []*Vote    `json:"votes"`
		Comments   []*Comment `json:"comments"`
		Created    time.Time  `json:"created"`
		UpVotePerc int        `json:"upvotePercentage"`
		ID         string     `json:"id"`
	}
	Vote struct {
		User string `json:"user"`
		Vote int    `json:"vote"`
	}
)
