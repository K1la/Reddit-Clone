package models

type (
	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	UserRepo interface {
		Create(user *User) (*User, error)
	}
)
