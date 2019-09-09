package models

import "time"

type SensorPoint struct {
	Time     time.Time `json:"t"`
	Co2Value float64   `json:"v"`
}

func (p *SensorPoint) CopyTo(other *SensorPoint) *SensorPoint {
	if p == nil {
		return nil
	}

	other.Time = p.Time
	other.Co2Value = p.Co2Value

	return other
}
