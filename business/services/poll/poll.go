package poll

import (
	"fmt"
	"goairmon/business/data/context"
	"goairmon/business/data/models"
	"goairmon/business/hardware"
	"sync"
	"time"

	"github.com/labstack/echo"
)

func NewPollService(cfg *Config, dbContext context.DbContext) *PollService {
	sensorCfg := &hardware.Co2SensorCfg{
		ReadDelayMillis:      1000,
		BaselineDelaySeconds: 60,
		Logger:               cfg.Logger,
	}
	co2Sensor := hardware.NewPiCo2Sensor(sensorCfg, dbContext)

	return &PollService{
		cfg:       cfg,
		co2Sensor: co2Sensor,
		stopChan:  nil,
		dbContext: dbContext,
	}
}

type PollService struct {
	dbContext context.DbContext
	stopChan  chan int
	lock      sync.Mutex
	cfg       *Config
	co2Sensor *hardware.Co2Sensor
}

type Config struct {
	PollDelayMillis int
	Logger          echo.Logger
}

func (p *PollService) Start() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.stopChan != nil {
		return fmt.Errorf("service already started")
	}

	if err := p.co2Sensor.Start(); err != nil {
		return err
	}

	p.stopChan = make(chan int)

	ticker := time.NewTicker(time.Millisecond * time.Duration(p.cfg.PollDelayMillis))
	go p.pollRoutine(ticker)

	return nil
}

func (p *PollService) Stop() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.stopChan == nil {
		return fmt.Errorf("service already stopped")
	}

	select {
	case p.stopChan <- 0:
		break
	case <-time.After(time.Millisecond * 100):
		break
	}

	if err := p.co2Sensor.Close(); err != nil {
		return err
	}

	p.stopChan = nil

	return nil
}

func (p *PollService) pollRoutine(pollTicker *time.Ticker) {
	defer pollTicker.Stop()

	for {
		select {
		case <-p.stopChan:
			return
		case <-pollTicker.C:
			if err := p.takePoll(); err != nil {
				p.cfg.Logger.Error("failed to poll sensor", err)
			}
		}
	}
}

func (p *PollService) takePoll() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.dbContext.PushSensorPoint(&models.SensorPoint{Time: time.Now(), Co2Value: float64(p.co2Sensor.ECO2)})
	if err != nil {
		return err
	}

	return p.dbContext.Save()
}
