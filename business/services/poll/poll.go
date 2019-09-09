package poll

import (
	"fmt"
	"goairmon/business/data/context"
	"goairmon/business/data/models"
	"goairmon/business/hardware"
	"sync"
	"time"
)

func NewPollService(cfg *Config) *PollService {
	sensorCfg := &hardware.Co2SensorCfg{}

	var co2Sensor hardware.Co2Sensor
	if cfg.PiSensor {
		co2Sensor = hardware.NewPiCo2Sensor(sensorCfg)
	} else {
		co2Sensor = hardware.NewFakeCo2Sensor(sensorCfg)
	}

	return &PollService{
		cfg:       cfg,
		co2Sensor: co2Sensor,
		stopChan:  nil,
	}
}

type PollService struct {
	stopChan  chan int
	lock      sync.Mutex
	cfg       *Config
	co2Sensor hardware.Co2Sensor
}

type Config struct {
	TickMillis int
	PiSensor   bool
	DbContext  context.DbContext
}

func (p *PollService) Start() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.stopChan != nil {
		return fmt.Errorf("service already started")
	}

	p.stopChan = make(chan int)

	go p.pollRoutine()

	if err := p.co2Sensor.On(); err != nil {
		return err
	}

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
		//
	case <-time.After(time.Millisecond * 100):
		//
	}

	if err := p.co2Sensor.Off(); err != nil {
		return err
	}

	return nil
}

func (p *PollService) pollRoutine() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(p.cfg.TickMillis))

	for {
		select {
		case <-p.stopChan:
			return
		case <-ticker.C:
			p.takePoll()
		}
	}
}

func (p *PollService) takePoll() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	val, err := p.co2Sensor.Read()
	if err != nil {
		return err
	}

	p.cfg.DbContext.PushSensorPoint(&models.SensorPoint{Time: time.Now(), Co2Value: val})

	return err
}
