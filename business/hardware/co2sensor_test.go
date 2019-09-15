package hardware

import (
	"goairmon/business/data/models"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStartStop(t *testing.T) {
	if runtime.GOARCH == "arm" {
		t.Log("skipping test in arm")
	}

	cfg := &Co2SensorCfg{
		ReadDelayMillis:      1000,
		BaselineDelaySeconds: 30,
	}

	dbContext := &_fakeDbContext{}
	dbContext.setBaselineClosure = func(eCO2 uint16, TVOC uint16) error {
		if eCO2 != 23 {
			t.Error("unexpected eco2 baseline value", 23, eCO2)
		}

		if TVOC != 42 {
			t.Error("unexpected TVOC value", TVOC)
		}

		return nil
	}

	dbContext.getBaselineClosure = func() (eCO2 uint16, TVOC uint16, err error) {
		return 23, 42, nil
	}

	co2Sensor := NewPiCo2Sensor(cfg, dbContext)

	co2Sensor.sgp30.(*fakeSgp30).actionDelayMillis = 10
	startTime := time.Now()

	if err := co2Sensor.Start(); err != nil {
		t.Error("unexpected err", err)
	}

	if err := co2Sensor.Start(); err == nil {
		t.Error("expected error")
	}

	if err := co2Sensor.Close(); err != nil {
		t.Error("unexpected err", err)
	}

	if err := co2Sensor.Close(); err == nil {
		t.Error("expected  err")
	}

	passedTime := time.Since(startTime)
	if passedTime < time.Millisecond*20 {
		t.Error("should have been delayed by action time")
	}
}

func TestLoopRoutine(t *testing.T) {
	if runtime.GOARCH == "arm" {
		t.Log("skipping test in arm")
	}

	cfg := &Co2SensorCfg{
		ReadDelayMillis:      1000,
		BaselineDelaySeconds: 30,
	}

	var eCO2Baseline uint16
	var tVOCBaseline uint16

	dbContext := &_fakeDbContext{}
	dbContext.setBaselineClosure = func(eCO2 uint16, TVOC uint16) error {
		eCO2Baseline = eCO2
		tVOCBaseline = TVOC

		return nil
	}

	dbContext.getBaselineClosure = func() (eCO2 uint16, TVOC uint16, err error) {
		return 23, 42, nil
	}

	co2Sensor := NewPiCo2Sensor(cfg, dbContext)

	readTicker := time.NewTicker(time.Hour)
	baselineTicker := time.NewTicker(time.Hour)
	defer func() {
		readTicker.Stop()
		baselineTicker.Stop()
	}()

	readChan := make(chan time.Time)
	readTicker.C = readChan
	baselineChan := make(chan time.Time)
	baselineTicker.C = baselineChan
	co2Sensor.stopChan = make(chan int)

	fakeSgp30 := co2Sensor.sgp30.(*fakeSgp30)
	fakeSgp30.staticECO2 = 1
	fakeSgp30.staticTVOC = 2

	go co2Sensor.loopRoutine(readTicker, baselineTicker)

	if co2Sensor.ECO2 != 0 || co2Sensor.TVOC != 0 {
		t.Error("unexpected start values")
	}

	readChan <- time.Now()

	if co2Sensor.ECO2 != 1 {
		t.Error("unexpected eco2 value", 1, co2Sensor.ECO2)
	}

	if co2Sensor.TVOC != 2 {
		t.Error("unexpected tvoc value", 2, co2Sensor.TVOC)
	}

	baselineChan <- time.Now()

	for range time.Tick(time.Millisecond) {
		if eCO2Baseline != 0 {
			break
		}
	}

	if eCO2Baseline != 1 {
		t.Error("unexpected ECO2 baseline", 1, eCO2Baseline)
	}
	if tVOCBaseline != 2 {
		t.Error("unexpected TVOC baseline", 2, tVOCBaseline)
	}
}

type _fakeDbContext struct {
	setBaselineClosure func(eCO2 uint16, TVOC uint16) error
	getBaselineClosure func() (eCO2 uint16, TVOC uint16, err error)
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
	panic("not implemented")
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
