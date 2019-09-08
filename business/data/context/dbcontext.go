package context

import (
	"goairmon/business/data/models"

	"github.com/google/uuid"
)

type DbContext interface {
	Close() error
	CreateOrUpdateUser(user *models.User) error
	FindUser(id uuid.UUID) (*models.User, error)
	FindUserByName(username string) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}
