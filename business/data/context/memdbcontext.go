package context

import (
	"encoding/json"
	"fmt"
	"goairmon/business/data/models"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

func NewMemDbContext(cfg *MemDbConfig) DbContext {
	ctx := &memDbContext{
		cfg: cfg,
	}

	if err := ctx.load(); err != nil {
		ctx.users = make(map[uuid.UUID]*models.User)
	}

	return ctx
}

type MemDbConfig struct {
	StoragePath string
}

type memDbContext struct {
	cfg   *MemDbConfig
	users map[uuid.UUID]*models.User
}

func (m *memDbContext) Close() error {
	return m.save()
}

func (m *memDbContext) CreateOrUpdateUser(user *models.User) error {
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
	existing, ok := m.users[id]
	if ok {
		user := existing.CopyTo(&models.User{})
		return user, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (m *memDbContext) FindUserByName(userName string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == userName {
			return user, nil
		}
	}

	return nil, fmt.Errorf("failed to find user")
}

func (m *memDbContext) DeleteUser(id uuid.UUID) error {
	_, ok := m.users[id]
	if ok {
		delete(m.users, id)
		return nil
	}

	return fmt.Errorf("id not found")
}

func (m *memDbContext) load() error {
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

func (m *memDbContext) save() error {
	raw, err := json.Marshal(m.users)

	if err != nil {
		return fmt.Errorf("failed to marshal user storage: %s", err)
	}

	err = ioutil.WriteFile(m.userFile(), raw, 0644)
	if err != nil {
		return fmt.Errorf("failed to save user storage: %s", err)
	}

	return nil
}

func (m *memDbContext) userFile() string {
	os.MkdirAll(m.cfg.StoragePath, 0700)

	return m.cfg.StoragePath + "/goairmon_users.json"
}
