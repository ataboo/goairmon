package hardware

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewFakeSgp30Sensor() SGP30 {
	min := uint16(1)
	max := uint16(1000)
	last := float64((max-min)/2) + float64(min)

	return &fakeSgp30{
		lastReading: last,
		min:         min,
		max:         max,
		variance:    0.1,
	}
}

type fakeSgp30 struct {
	lastReading       float64
	min               uint16
	max               uint16
	actionDelayMillis int
	variance          float64
	staticECO2        uint16
	staticTVOC        uint16
}

// Init() error
// 	Close() error
// 	Measure() (eCO2 uint16, TVOC uint16, err error)
// 	GetBaseline() (eCO2 uint16, TVOC uint16, err error)
// 	SetBaseline(eCO2 uint16, TVOC uint16) error

func (s *fakeSgp30) Init() error {
	if s.actionDelayMillis > 0 {
		time.Sleep(time.Duration(s.actionDelayMillis) * time.Millisecond)
	}

	return nil
}

func (s *fakeSgp30) GetBaseline() (eCO2 uint16, TVOC uint16, err error) {
	if s.actionDelayMillis > 0 {
		time.Sleep(time.Duration(s.actionDelayMillis) * time.Millisecond)
	}

	if s.staticECO2 != 0 || s.staticTVOC != 0 {
		return s.staticECO2, s.staticTVOC, nil
	}

	return 12, 34, nil
}

func (s *fakeSgp30) SetBaseline(eCO2 uint16, TVOC uint16) error {
	if s.actionDelayMillis > 0 {
		time.Sleep(time.Duration(s.actionDelayMillis) * time.Millisecond)
	}

	return nil
}

func (s *fakeSgp30) Measure() (eCO2 uint16, TVOC uint16, err error) {
	if s.actionDelayMillis > 0 {
		time.Sleep(time.Duration(s.actionDelayMillis) * time.Millisecond)
	}

	if s.staticECO2 != 0 || s.staticTVOC != 0 {
		return s.staticECO2, s.staticTVOC, nil
	}

	delta := float64(s.max - s.min)
	newVal := s.lastReading + (rand.Float64()*2.0-1.0)*delta*s.variance
	newVal = math.Min(math.Max(newVal, float64(s.min)), float64(s.max))
	s.lastReading = newVal

	return uint16(newVal), uint16(newVal), nil
}

func (s *fakeSgp30) Close() error {
	if s.actionDelayMillis > 0 {
		time.Sleep(time.Duration(s.actionDelayMillis) * time.Millisecond)
	}

	return nil
}
