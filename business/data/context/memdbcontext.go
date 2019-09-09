package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"goairmon/business/data/models"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func NewMemDbContext(cfg *MemDbConfig) DbContext {
	ctx := &memDbContext{
		cfg:          cfg,
		sensorPoints: NewSensorPointStack(cfg.SensorPointCount),
	}

	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if err := ctx.loadUsers(); err != nil {
		ctx.users = make(map[uuid.UUID]*models.User)
	}

	if err := ctx.loadPoints(); err != nil {
		ctx.sensorPoints.Clear()
	}

	return ctx
}

type MemDbConfig struct {
	StoragePath      string
	SensorPointCount int
	EncodeReadible   bool
}

type memDbContext struct {
	cfg          *MemDbConfig
	users        map[uuid.UUID]*models.User
	sensorPoints PointStack
	lock         sync.Mutex
}

func (m *memDbContext) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	errs := make([]string, 0)

	if err := m.savePoints(); err != nil {
		errs = append(errs, err.Error())
	}

	if err := m.saveUsers(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func (m *memDbContext) CreateOrUpdateUser(user *models.User) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if user.ID != uuid.Nil {
		existing, ok := m.users[user.ID]
		if ok {
			user.CopyTo(existing)
			return nil
		}
	}

	user.ID = uuid.New()
	m.users[user.ID] = user.CopyTo(&models.User{})

	return nil
}

func (m *memDbContext) FindUser(id uuid.UUID) (*models.User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	existing, ok := m.users[id]
	if ok {
		user := existing.CopyTo(&models.User{})
		return user, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (m *memDbContext) FindUserByName(userName string) (*models.User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, user := range m.users {
		if user.Username == userName {
			return user, nil
		}
	}

	return nil, fmt.Errorf("failed to find user")
}

func (m *memDbContext) DeleteUser(id uuid.UUID) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.users[id]
	if ok {
		delete(m.users, id)
		return nil
	}

	return fmt.Errorf("id not found")
}

func (m *memDbContext) loadUsers() error {
	raw, err := ioutil.ReadFile(m.userFile())
	if err != nil {
		return fmt.Errorf("failed to read user storage: %s", err)
	}
	err = json.Unmarshal(raw, &m.users)

	if err != nil {
		return fmt.Errorf("failed to decode user storage: %s", err)
	}

	return nil
}

func (m *memDbContext) loadPoints() error {
	raw, err := ioutil.ReadFile(m.pointFile())
	if err != nil {
		return fmt.Errorf("failed to read point storage: %s", err)
	}

	if err := m.sensorPoints.Decode(raw); err != nil {
		return fmt.Errorf("failed to decode point storage: %s", err)
	}

	return nil
}

func (m *memDbContext) saveUsers() error {
	raw, err := json.Marshal(m.users)
	if err != nil {
		return fmt.Errorf("failed to marshal user storage: %s", err)
	}

	os.MkdirAll(m.cfg.StoragePath, 0700)
	if err := ioutil.WriteFile(m.userFile(), raw, 0644); err != nil {
		return fmt.Errorf("failed to save user storage: %s", err)
	}

	return nil
}

func (m *memDbContext) savePoints() error {
	raw, err := m.sensorPoints.Encode()
	if err != nil {
		return fmt.Errorf("failed to marshal sensor points: %s", err)
	}

	os.MkdirAll(m.cfg.StoragePath, 0700)
	if err := ioutil.WriteFile(m.pointFile(), raw, 0644); err != nil {
		return fmt.Errorf("failed to write sensor points: %s", err)
	}

	return nil
}

func (m *memDbContext) userFile() string {
	return m.cfg.StoragePath + "/goairmon_users.json"
}

func (m *memDbContext) pointFile() string {
	return m.cfg.StoragePath + "/goairmon_points.json"
}

func (m *memDbContext) PushSensorPoint(point *models.SensorPoint) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.sensorPoints.Push(point)

	return nil
}

func (m *memDbContext) GetSensorPoints(count int) ([]*models.SensorPoint, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.sensorPoints.PeakNLatest(count)
}
