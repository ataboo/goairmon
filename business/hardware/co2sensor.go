package hardware

import "fmt"

type Co2Sensor interface {
	IsOn() bool
	On() error
	Off() error
	Read() (float64, error)
	Close()
}

type Co2SensorCfg struct {
	PinNumber uint
}

func NewPiCo2Sensor(cfg *Co2SensorCfg) Co2Sensor {
	sensor := &piCo2Sensor{
		cfg: cfg,
	}

	if err := sensor.init(); err != nil {
		fmt.Println("failed to init piCo2Sensor")
		return nil
	}

	return sensor
}

type piCo2Sensor struct {
	cfg    *Co2SensorCfg
	active bool
}

func (s *piCo2Sensor) init() error {
	//
	return nil
}

func (s *piCo2Sensor) On() error {
	s.active = true

	return nil
}

func (s *piCo2Sensor) Off() error {
	s.active = false

	return nil
}

func (s *piCo2Sensor) IsOn() bool {
	return s.active
}

func (s *piCo2Sensor) Read() (float64, error) {
	//

	return 0, nil
}

func (s *piCo2Sensor) Close() {
	//
}
