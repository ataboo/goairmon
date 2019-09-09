package hardware

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewFakeCo2Sensor(cfg *Co2SensorCfg) Co2Sensor {
	min := 0.0
	max := 100.0

	sensor := &fakeCo2Sensor{
		cfg:         cfg,
		active:      false,
		lastReading: (max-min)/2.0 + min,
		min:         min,
		max:         max,
		variance:    0.1,
	}

	return sensor
}

type fakeCo2Sensor struct {
	cfg         *Co2SensorCfg
	active      bool
	lastReading float64
	min         float64
	max         float64
	variance    float64
}

func (s *fakeCo2Sensor) On() error {
	s.active = true

	return nil
}

func (s *fakeCo2Sensor) Off() error {
	s.active = false

	return nil
}

func (s *fakeCo2Sensor) IsOn() bool {
	return s.active
}

func (s *fakeCo2Sensor) Read() (float64, error) {
	if !s.active {
		return -1, fmt.Errorf("sensor not active")
	}

	newVal := s.lastReading + (rand.Float64()*2.0-1.0)*(s.max-s.min)*s.variance
	newVal = math.Min(math.Max(newVal, s.min), s.max)
	s.lastReading = newVal

	return newVal, nil
}

func (s *fakeCo2Sensor) Close() {
	//
}
