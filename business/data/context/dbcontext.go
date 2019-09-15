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
	PushSensorPoint(point *models.SensorPoint) error
	GetSensorPoints(count int) ([]*models.SensorPoint, error)
	GetSensorBaseline() (eCO2 uint16, TVOC uint16, err error)
	SetSensorBaseline(eCO2 uint16, TVOC uint16) error
	Save() error
}
