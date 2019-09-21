package hardware

import (
	"fmt"
	"goairmon/business/data/context"
	"runtime"
	"sync"
	"time"

	"github.com/ataboo/sgp30go/sensor"
	"github.com/labstack/echo"
)

type SGP30 interface {
	Init() error
	Close() error
	Measure() (eCO2 uint16, TVOC uint16, err error)
	GetBaseline() (eCO2 uint16, TVOC uint16, err error)
	SetBaseline(eCO2 uint16, TVOC uint16) error
}

type Co2SensorCfg struct {
	ReadDelayMillis      int
	BaselineDelaySeconds int
	Logger               echo.Logger
}

func NewPiCo2Sensor(cfg *Co2SensorCfg, dbContext context.DbContext) *Co2Sensor {
	co2Sensor := &Co2Sensor{
		cfg:       cfg,
		dbContext: dbContext,
	}

	if runtime.GOARCH == "arm" {
		cfg.Logger.Info("Detected arm, starting i2c sensor")
		sensorCfg := sensor.DefaultConfig()
		sensorCfg.Logger = NewLoggingAdaptor(cfg.Logger)
		co2Sensor.sgp30 = sensor.NewSensor(sensorCfg)
	} else {
		cfg.Logger.Info("Detected non-arm, starting fake sensor values")
		co2Sensor.sgp30 = newFakeSgp30Sensor()
	}

	return co2Sensor
}

type Co2Sensor struct {
	cfg       *Co2SensorCfg
	sgp30     SGP30
	stopChan  chan int
	lock      sync.Mutex
	ECO2      uint16
	TVOC      uint16
	dbContext context.DbContext
}

func (s *Co2Sensor) Start() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.stopChan != nil {
		return fmt.Errorf("sensor already started")
	}

	if err := s.sgp30.Init(); err != nil {
		return err
	}

	s.applySavedSensorBaseline()

	s.stopChan = make(chan int)
	readTicker := time.NewTicker(time.Millisecond * time.Duration(s.cfg.ReadDelayMillis))
	baselineTicker := time.NewTicker(time.Second * time.Duration(s.cfg.BaselineDelaySeconds))
	go s.loopRoutine(readTicker, baselineTicker)

	return nil
}

func (s *Co2Sensor) applySavedSensorBaseline() {
	eCO2, TVOC, err := s.dbContext.GetSensorBaseline()
	if err != nil {
		s.cfg.Logger.Error("failed to load saved sensor baseline", err)
	}

	if eCO2 > 0 && TVOC > 0 {
		if err := s.sgp30.SetBaseline(eCO2, TVOC); err != nil {
			s.cfg.Logger.Error("failed to set sgp30 baseline", err)
		}
	}
}

func (s *Co2Sensor) loopRoutine(readTicker *time.Ticker, baseLineTicker *time.Ticker) {
	defer func() {
		readTicker.Stop()
		baseLineTicker.Stop()
	}()

	for {
		select {
		case <-s.stopChan:
			return
		case <-readTicker.C:
			eCO2, TVOC, err := s.sgp30.Measure()
			if err != nil {
				s.cfg.Logger.Error("failed to measure", err)
			}
			s.ECO2 = eCO2
			s.TVOC = TVOC
		case <-baseLineTicker.C:
			eCO2, TVOC, err := s.sgp30.GetBaseline()
			if err != nil {
				s.cfg.Logger.Error("failed to get baseline", err)
			} else {
				if err := s.dbContext.SetSensorBaseline(eCO2, TVOC); err != nil {
					s.cfg.Logger.Error(err)
				}
			}
		}
	}
}

func (s *Co2Sensor) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.stopChan == nil {
		return fmt.Errorf("sensor already stopped")
	}

	select {
	case s.stopChan <- 0:
		break
	case <-time.After(time.Millisecond * 100):
		s.cfg.Logger.Error("co2sensor stop signal timed out")
		break
	}

	s.stopChan = nil

	return s.sgp30.Close()
}
