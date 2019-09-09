package context

import (
	"encoding/json"
	"fmt"
	"goairmon/business/data/models"
)

type PointStack interface {
	Push(point *models.SensorPoint)
	Peak(idx int) *models.SensorPoint
	PeakNLatest(count int) ([]*models.SensorPoint, error)
	Encode() ([]byte, error)
	Decode(raw []byte) error
	Pop() *models.SensorPoint
	Size() int
	Clear()
}

func NewSensorPointStack(size int) PointStack {
	stack := &sensorPointStack{
		size: size,
	}

	stack.Clear()

	return stack
}

type sensorPointStack struct {
	Index  int                   `json:"index"`
	Values []*models.SensorPoint `json:"values"`
	size   int
}

func (s *sensorPointStack) Push(point *models.SensorPoint) {
	s.Index = s.normalizedIdx(s.Index + 1)
	s.Values[s.Index] = point
}

func (s *sensorPointStack) Pop() *models.SensorPoint {
	val := s.Values[s.Index]
	s.Values[s.Index] = nil
	s.Index = s.normalizedIdx(s.Index - 1)

	return val
}

func (s *sensorPointStack) Peak(idx int) *models.SensorPoint {
	return s.Values[s.normalizedIdx(s.Index-idx)]
}

func (s *sensorPointStack) PeakNLatest(count int) ([]*models.SensorPoint, error) {
	if count < 1 {
		count = s.size
	} else if count > s.size {
		return nil, fmt.Errorf("count %d larger than stack size %d", count, s.size)
	}

	out := make([]*models.SensorPoint, count)

	for i := 0; i < count; i++ {
		out[i] = s.Peak(i).CopyTo(&models.SensorPoint{})
	}

	return out, nil
}

func (s *sensorPointStack) Size() int {
	return s.size
}

func (s *sensorPointStack) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *sensorPointStack) Decode(raw []byte) error {
	err := json.Unmarshal(raw, s)
	s.size = len(s.Values)

	return err
}

func (s *sensorPointStack) Clear() {
	s.Index = s.size - 1
	s.Values = make([]*models.SensorPoint, s.size)
}

func (s *sensorPointStack) normalizedIdx(idx int) int {
	mod := idx % s.Size()
	if mod < 0 {
		mod += s.Size()
	}

	return mod
}
