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
	"github.com/labstack/echo"
)

func NewMemDbContext(cfg *MemDbConfig) DbContext {
	if cfg.SensorPointCount == 0 {
		cfg.SensorPointCount = 48 * 60
	}

	ctx := &memDbContext{
		cfg:          cfg,
		sensorPoints: NewSensorPointStack(cfg.SensorPointCount),
	}

	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if err := ctx.loadStoredConfig(); err != nil {
		ctx.storedConfig = &StoredConfig{
			Users: make(map[uuid.UUID]*models.User),
		}
	}

	if err := ctx.loadPoints(); err != nil {
		ctx.sensorPoints.Clear()
	}

	return ctx
}

type StoredConfig struct {
	ECO2Baseline uint16 `json:"eco2"`
	TVOCBaseline uint16 `json:"tvoc"`
	Users        map[uuid.UUID]*models.User
}

type MemDbConfig struct {
	StoragePath      string
	SensorPointCount int
	EncodeReadible   bool
	Logger           echo.Logger
}

type memDbContext struct {
	cfg          *MemDbConfig
	sensorPoints PointStack
	storedConfig *StoredConfig
	lock         sync.Mutex
}

func (m *memDbContext) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	errs := make([]string, 0)

	if err := m.savePoints(); err != nil {
		errs = append(errs, err.Error())
	}

	if err := m.saveStoredConfig(); err != nil {
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
		existing, ok := m.storedConfig.Users[user.ID]
		if ok {
			user.CopyTo(existing)
			return nil
		}
	}

	user.ID = uuid.New()
	m.storedConfig.Users[user.ID] = user.CopyTo(&models.User{})

	return nil
}

func (m *memDbContext) FindUser(id uuid.UUID) (*models.User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	existing, ok := m.storedConfig.Users[id]
	if ok {
		user := existing.CopyTo(&models.User{})
		return user, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (m *memDbContext) FindUserByName(userName string) (*models.User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, user := range m.storedConfig.Users {
		if user.Username == userName {
			return user, nil
		}
	}

	return nil, fmt.Errorf("failed to find user")
}

func (m *memDbContext) DeleteUser(id uuid.UUID) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.storedConfig.Users[id]
	if ok {
		delete(m.storedConfig.Users, id)
		return nil
	}

	return fmt.Errorf("id not found")
}

func (m *memDbContext) loadPoints() error {
	raw, err := ioutil.ReadFile(m.pointFile())
	if err != nil {
		return fmt.Errorf("failed to read point storage: %s", err)
	}

	if err := json.Unmarshal(raw, m.sensorPoints); err != nil {
		return fmt.Errorf("failed to decode point storage: %s", err)
	}

	if m.sensorPoints.Size() != m.cfg.SensorPointCount {
		m.sensorPoints.Resize(m.cfg.SensorPointCount)
	}

	return nil
}

func (m *memDbContext) loadStoredConfig() error {
	raw, err := ioutil.ReadFile(m.configFile())
	if err != nil {
		return fmt.Errorf("failed to read stored config: %s", err)
	}

	if err := json.Unmarshal(raw, &m.storedConfig); err != nil {
		return fmt.Errorf("failed to decode stored config: %s", err)
	}

	return nil
}

func (m *memDbContext) saveStoredConfig() error {
	raw, err := json.Marshal(m.storedConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal stored config: %s", err)
	}

	os.MkdirAll(m.cfg.StoragePath, 0700)
	if err := ioutil.WriteFile(m.configFile(), raw, 0644); err != nil {
		return fmt.Errorf("failed to save user storage: %s", err)
	}

	return nil
}

func (m *memDbContext) savePoints() error {
	raw, err := json.Marshal(m.sensorPoints)
	if err != nil {
		return fmt.Errorf("failed to marshal sensor points: %s", err)
	}

	os.MkdirAll(m.cfg.StoragePath, 0700)
	if err := ioutil.WriteFile(m.pointFile(), raw, 0644); err != nil {
		return fmt.Errorf("failed to write sensor points: %s", err)
	}

	return nil
}

func (m *memDbContext) configFile() string {
	return m.cfg.StoragePath + "/goairmon_config.json"
}

func (m *memDbContext) pointFile() string {
	return m.cfg.StoragePath + "/goairmon_points.json"
}

func (m *memDbContext) PushSensorPoint(point *models.SensorPoint) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	lastPoint := m.sensorPoints.Peak(0)
	if lastPoint != nil && lastPoint.Time.Day() != point.Time.Day() {
		if err := m.archiveLastDay(); err != nil {
			m.cfg.Logger.Error(err)
		}
	}

	m.sensorPoints.Push(point)

	return nil
}

func (m *memDbContext) archiveLastDay() error {
	daysPoints, err := m.sensorPoints.PeakNLatest(60 * 24)
	if err != nil {
		return err
	}

	if len(daysPoints) == 0 || daysPoints[0] == nil {
		return fmt.Errorf("not enough points saved")
	}

	year := daysPoints[0].Time.Year()
	day := daysPoints[0].Time.YearDay()

	for i := len(daysPoints) - 1; i >= 0; i-- {
		point := daysPoints[i]
		if point == nil || point.Time.YearDay() != day || point.Time.Year() != year {
			continue
		}

		daysPoints = daysPoints[:i]
		break
	}

	if len(daysPoints) == 0 {
		return fmt.Errorf("not enough valid points")
	}

	encoded, err := json.Marshal(daysPoints)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("/archive_%d_%02d_%02d.json", daysPoints[0].Time.Year(), daysPoints[0].Time.Month(), daysPoints[0].Time.Day())
	os.MkdirAll(m.cfg.StoragePath, 0700)
	if err := ioutil.WriteFile(m.cfg.StoragePath+fileName, encoded, 0644); err != nil {
		return fmt.Errorf("failed to write archive: %s", err)
	}

	return nil
}

func (m *memDbContext) GetSensorPoints(count int) ([]*models.SensorPoint, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	out := make([]*models.SensorPoint, 0)

	points, err := m.sensorPoints.PeakNLatest(count)
	if err != nil {
		return out, err
	}

	for _, p := range points {
		if p != nil {
			out = append(out, p)
		}
	}

	return out, nil
}

func (m *memDbContext) ClearSensorPoints() error {
	m.sensorPoints.Clear()

	return m.savePoints()
}

func (m *memDbContext) GetSensorBaseline() (eCO2 uint16, TVOC uint16, err error) {
	eCO2 = m.storedConfig.ECO2Baseline
	TVOC = m.storedConfig.TVOCBaseline
	if eCO2 == 0 || TVOC == 0 {
		err = fmt.Errorf("baseline not set")
	}
	return eCO2, TVOC, err
}

func (m *memDbContext) SetSensorBaseline(eCO2 uint16, TVOC uint16) error {
	m.storedConfig.ECO2Baseline = eCO2
	m.storedConfig.TVOCBaseline = TVOC

	return nil
}

func (m *memDbContext) Save() error {
	m.saveStoredConfig()
	return m.savePoints()
}
