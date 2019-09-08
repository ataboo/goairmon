package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `col:"id"`
	Username     string    `col:"username"`
	PasswordHash string    `col:"passwordhash"`
	LastLogin    time.Time `col:"lastlogin"`
}

func (u *User) CopyTo(other *User) *User {
	other.ID = u.ID
	other.Username = u.Username
	other.PasswordHash = u.PasswordHash
	other.LastLogin = u.LastLogin

	return other
}
