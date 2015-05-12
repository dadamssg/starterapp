package user

import (
	"time"
)

type User struct {
	Id                string
	CreatedAt         time.Time
	Username          string
	Email             string
	Password          string
	Enabled           bool
	ConfirmationToken string
}

type UserRepository interface {
	ById(id string) (*User, error)
	ByEmail(email string) (*User, error)
	ByUsername(username string) (*User, error)
	Add(user *User) error
}
