package models

type (
	Session struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	SessionRepo interface {
		Create(*Session) (*Session, error)
	}
)
