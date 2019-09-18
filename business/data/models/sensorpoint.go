package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type SensorPoint struct {
	Time     time.Time
	Co2Value float64
}

type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).Unix())), nil
}

func (t *JsonTime) UnmarshalJSON(raw []byte) error {
	stamp, err := strconv.ParseInt(string(raw), 0, 64)
	if err != nil {
		return err
	}

	parsedTime := time.Unix(stamp, 0)

	*t = JsonTime(parsedTime)

	return err
}

func (p *SensorPoint) MarshalJSON() ([]byte, error) {
	jsonStruct := struct {
		JTime    JsonTime `json:"t"`
		Co2Value float64  `json:"v"`
	}{
		JsonTime(p.Time),
		p.Co2Value,
	}

	return json.Marshal(jsonStruct)
}

func (p *SensorPoint) UnmarshalJSON(raw []byte) error {
	jsonStruct := struct {
		JTime    JsonTime `json:"t"`
		Co2Value float64  `json:"v"`
	}{}

	err := json.Unmarshal(raw, &jsonStruct)
	if err != nil {
		return err
	}

	p.Time = time.Time(jsonStruct.JTime)
	p.Co2Value = jsonStruct.Co2Value

	return nil
}

func (p *SensorPoint) CopyTo(other *SensorPoint) *SensorPoint {
	if p == nil {
		return nil
	}

	other.Time = p.Time
	other.Co2Value = p.Co2Value

	return other
}
