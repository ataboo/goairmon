package poll

import (
	"goairmon/business/data/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func TestPollStart(t *testing.T) {
	cfg := &Config{
		PollDelayMillis: 1000 * 60,
		Logger:          echo.New().Logger,
	}
	ctx := &_fakeDbContext{
		setBaselineClosure: func(eCO2 uint16, tVOC uint16) error {
			return nil
		},
		getBaselineClosure: func() (eCO2 uint16, tVOC uint16, err error) {
			return 1, 2, nil
		},
	}

	poll := NewPollService(cfg, ctx)

	if err := poll.Start(); err != nil {
		t.Error(err)
	}

	if err := poll.Start(); err == nil {
		t.Error("expected error")
	}

	if err := poll.Stop(); err != nil {
		t.Error(err)
	}

	if err := poll.Stop(); err == nil {
		t.Error("expected error")
	}
}

func TestPollRoutine(t *testing.T) {
	sensorPoints := make([]*models.SensorPoint, 0)
	ctx := &_fakeDbContext{
		setBaselineClosure: func(eCO2 uint16, tVOC uint16) error {
			return nil
		},
		getBaselineClosure: func() (eCO2 uint16, tVOC uint16, err error) {
			return 23, 2, nil
		},
		sensorPointClosure: func(point *models.SensorPoint) error {
			sensorPoints = append(sensorPoints, point)

			return nil
		},
	}

	poll := NewPollService(&Config{Logger: echo.New().Logger}, ctx)

	poll.stopChan = make(chan int)

	ticker := time.NewTicker(time.Second)
	tickChan := make(chan time.Time)
	ticker.C = tickChan
	poll.co2Sensor.ECO2 = 23
	go poll.pollRoutine(ticker)

	if len(sensorPoints) != 0 {
		t.Error("unexpected sensor point count", 0, len(sensorPoints))
	}

	tickChan <- time.Now()

	if len(sensorPoints) != 1 {
		t.Error("unexpected sensor point count", 1, len(sensorPoints))
	}

	if sensorPoints[0].Co2Value != 23 {
		t.Error("unexpected co2 value", 23, sensorPoints[0].Co2Value)
	}

	poll.stopChan <- 0
}

type _fakeDbContext struct {
	setBaselineClosure func(eCO2 uint16, TVOC uint16) error
	getBaselineClosure func() (eCO2 uint16, TVOC uint16, err error)
	sensorPointClosure func(point *models.SensorPoint) error
}

func (f *_fakeDbContext) Close() error {
	panic("not implemented")
}

func (f *_fakeDbContext) CreateOrUpdateUser(user *models.User) error {
	panic("not implemented")
}

func (f *_fakeDbContext) FindUser(id uuid.UUID) (*models.User, error) {
	panic("not implemented")
}

func (f *_fakeDbContext) FindUserByName(username string) (*models.User, error) {
	panic("not implemented")
}

func (f *_fakeDbContext) DeleteUser(id uuid.UUID) error {
	panic("not implemented")
}

func (f *_fakeDbContext) PushSensorPoint(point *models.SensorPoint) error {
	return f.sensorPointClosure(point)
}

func (f *_fakeDbContext) GetSensorPoints(count int) ([]*models.SensorPoint, error) {
	panic("not implemented")
}

func (f *_fakeDbContext) GetSensorBaseline() (eCO2 uint16, TVOC uint16, err error) {
	return f.getBaselineClosure()
}

func (f *_fakeDbContext) SetSensorBaseline(eCO2 uint16, TVOC uint16) error {
	return f.setBaselineClosure(eCO2, TVOC)
}

func (f *_fakeDbContext) Save() error {
	return nil
}

func (f *_fakeDbContext) ClearSensorPoints() error {
	panic("not implemented")
}
