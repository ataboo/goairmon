package context

import (
	"goairmon/business/data/models"
	"testing"
	"time"
)

func TestNormalizedIdx(t *testing.T) {
	rows := []struct {
		idx      int
		expected int
	}{
		{-9, 3},
		{-8, 0},
		{-7, 1},
		{-6, 2},
		{-5, 3},
		{-4, 0},
		{-3, 1},
		{-2, 2},
		{-1, 3},
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 0},
		{5, 1},
		{6, 2},
		{7, 3},
		{8, 0},
		{9, 1},
	}

	stack := sensorPointStack{
		Index:  3,
		Values: make([]*models.SensorPoint, 4),
		size:   4,
	}

	for _, row := range rows {
		normalized := stack.normalizedIdx(row.idx)
		if normalized != row.expected {
			t.Error("unexpected normalized value", row.idx, row.expected, normalized)
		}
	}
}

func TestNewSensorPointStack(t *testing.T) {
	stack := NewSensorPointStack(4)
	downcast := stack.(*sensorPointStack)

	points := []*models.SensorPoint{
		{Co2Value: 1.0},
		{Co2Value: 2.0},
		{Co2Value: 3.0},
		{Co2Value: 4.0},
	}

	if downcast.Index != 3 {
		t.Error("unexpected index", 3, downcast.Index)
	}

	for _, p := range points {
		stack.Push(p)
	}

	if downcast.Index != 3 {
		t.Error("unnexpected index", 3, downcast.Index)
	}

	if downcast.Values[0].Co2Value != 1.0 || stack.Peak(3).Co2Value != 1.0 || stack.Peak(7).Co2Value != 1.0 {
		t.Error("unexpected value", 1.0, downcast.Values[0].Co2Value)
	}
	if downcast.Values[1].Co2Value != 2.0 || stack.Peak(2).Co2Value != 2.0 || stack.Peak(6).Co2Value != 2.0 {
		t.Error("unexpected value", 2.0, downcast.Values[1].Co2Value)
	}
	if downcast.Values[2].Co2Value != 3.0 || stack.Peak(1).Co2Value != 3.0 || stack.Peak(5).Co2Value != 3.0 {
		t.Error("unexpected value", 3.0, downcast.Values[2].Co2Value)
	}
	if downcast.Values[3].Co2Value != 4.0 || stack.Peak(0).Co2Value != 4.0 || stack.Peak(4).Co2Value != 4.0 {
		t.Error("unexpected value", 4.0, downcast.Values[3].Co2Value)
	}

	stack.Push(&models.SensorPoint{Co2Value: 5.0})

	if downcast.Values[1].Co2Value != 2.0 || stack.Peak(3).Co2Value != 2.0 || stack.Peak(7).Co2Value != 2.0 {
		t.Error("unexpected value", 2.0, downcast.Values[1].Co2Value)
	}
	if downcast.Values[2].Co2Value != 3.0 || stack.Peak(2).Co2Value != 3.0 || stack.Peak(6).Co2Value != 3.0 {
		t.Error("unexpected value", 3.0, downcast.Values[2].Co2Value)
	}
	if downcast.Values[3].Co2Value != 4.0 || stack.Peak(1).Co2Value != 4.0 || stack.Peak(5).Co2Value != 4.0 {
		t.Error("unexpected value", 4.0, downcast.Values[3].Co2Value)
	}
	if downcast.Values[0].Co2Value != 5.0 || stack.Peak(0).Co2Value != 5.0 || stack.Peak(4).Co2Value != 5.0 {
		t.Error("unexpected value", 5.0, downcast.Values[0].Co2Value)
	}

	val := stack.Pop()

	if val.Co2Value != 5.0 {
		t.Error("unexpected value", 5.0, val.Co2Value)
	}

	if downcast.Index != 3 {

		t.Error("unexpected index", 3, downcast.Index)

	}

	if downcast.Values[0] != nil {
		t.Error("expected nil value", downcast.Values[0])
	}
}

func TestPeakNLatest(t *testing.T) {
	stack := NewSensorPointStack(4)
	// downcast := stack.(*sensorPointStack)

	points := []*models.SensorPoint{
		{Co2Value: 1.0, Time: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Co2Value: 2.0, Time: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Co2Value: 3.0, Time: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, p := range points {
		stack.Push(p)
	}

	result, err := stack.PeakNLatest(2)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 2 {
		t.Error("unexpected result length", 2, len(result))
	}

	if result[0].Co2Value != 3.0 || result[0].Time.Year() != 2012 {
		t.Error("unexpected result values", points[2], result[0])
	}

	if result[1].Co2Value != 2.0 || result[1].Time.Year() != 2011 {
		t.Error("unexpected result values", points[1], result[1])
	}

	result[0].Co2Value = 5.0
	if stack.Peak(0).Co2Value != 3.0 {
		t.Error("expected original value to be unchanged")
	}

	if _, err := stack.PeakNLatest(5); err == nil {
		t.Error("expected error")
	}

	result, err = stack.PeakNLatest(0)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 4 {
		t.Error("unexpected result length", 4, len(result))
	}
}
