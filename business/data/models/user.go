package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `col:"id"`
	Username     string    `col:"username"`
	PasswordHash []byte    `col:"passwordhash"`
	LastLogin    time.Time `col:"lastlogin"`
}

func (u *User) CopyTo(other *User) *User {
	other.ID = u.ID
	other.Username = u.Username
	other.PasswordHash = u.PasswordHash
	other.LastLogin = u.LastLogin

	return other
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)) == nil
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		u.PasswordHash = hashed
	}

	return err
}
