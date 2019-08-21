package models

import (
	"time"
)

type User struct {
	Id           int       `col:"id"`
	Username     string    `col:"username"`
	PasswordHash string    `col:"passwordhash"`
	LastLogin    time.Time `col:"lastlogin"`
}
