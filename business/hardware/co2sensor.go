package hardware

import (
	"goairmon/business/data/context"
	"runtime"

	"github.com/ataboo/spg30go/sensor"
	"github.com/op/go-logging"
)

type SGP30 interface {
	Init() error
	Close() error
	Measure() (eCO2 uint16, TVOC uint16, err error)
	GetBaseline() (eCO2 uint16, TVOC uint16, err error)
	SetBaseline(eCO2 uint16, TVOC uint16) error
}

type Co2SensorCfg struct {
	ReadDelayMillis       int
	CalibrateDelaySeconds int
	dbContext             context.DbContext
}

func NewPiCo2Sensor(cfg *Co2SensorCfg) *Co2Sensor {
	var sgp30 *sensor.SGP30Sensor

	if runtime.GOARCH == "arm" {
		sensorCfg := sensor.DefaultConfig()
		sensorCfg.Logger = logging.MustGetLogger("goairmon-spg30")
		sgp30 = sensor.NewSensor(sensorCfg)
	} else {
		sgp30 = fakeSgp30Sensor
	}

	sensor := &Co2Sensor{
		cfg:   cfg,
		sgp30: sgp30,
	}

	return sensor
}

type Co2Sensor struct {
	cfg   *Co2SensorCfg
	sgp30 SGP30
}

func (s *Co2Sensor) Start() error {
	if err := s.sgp30.Init(); err != nil {
		return err
	}

	return nil
}

func (s *Co2Sensor) Co2Value() (float64, error) {
	//

	return 0, nil
}

func (s *Co2Sensor) Close() {
	s.sgp30.Close()
}
