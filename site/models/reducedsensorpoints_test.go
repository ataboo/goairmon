package models

import (
	"goairmon/business/data/models"
	"testing"
	"time"
)

func TestReduce2Hours(t *testing.T) {
	startTime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

	rawPoints := []*models.SensorPoint{
		&models.SensorPoint{Time: startTime, Co2Value: 10},
		&models.SensorPoint{Time: startTime.Add(-time.Minute), Co2Value: 11},
		&models.SensorPoint{Time: startTime.Add(-time.Minute * time.Duration(2)), Co2Value: 12},
		&models.SensorPoint{Time: startTime.Add(-time.Minute * time.Duration(3)), Co2Value: 13},
		&models.SensorPoint{Time: startTime.Add(-time.Minute * time.Duration(4)), Co2Value: 14},
		&models.SensorPoint{Time: startTime.Add(-time.Minute * time.Duration(5)), Co2Value: 15},
		&models.SensorPoint{Time: startTime.Add(-time.Minute * time.Duration(6)), Co2Value: 16},
		&models.SensorPoint{Time: startTime.Add(-time.Minute * time.Duration(7)), Co2Value: 17},
	}

	reducedPoints := NewReducedSensorPoints(rawPoints, startTime)

	twoHours := reducedPoints.Last2Hours()

	if len(twoHours) != 120 {
		t.Error("unexpected count", 120, len(twoHours))
	}

	if !twoHours[0].Time.Equal(startTime) {
		t.Error("unexpected first time value", startTime, twoHours[0].Time)
	}

	if twoHours[0].Co2Value != 10 {
		t.Error("unexpected first co2 value", 10, twoHours[0].Co2Value)
	}

	if !twoHours[7].Time.Equal(startTime.Add(-time.Minute * time.Duration(7))) {
		t.Error("unexpected 7th item time")
	}

	if twoHours[7].Co2Value != 17.0 {
		t.Error("unexpected 7th item co2 value", 17.0, twoHours[7].Co2Value)
	}

	if twoHours[8].Co2Value != 400.0 {
		t.Error("unexpected 8th item co2 value", 400.0, twoHours[8].Co2Value)
	}

	mean := reducedPoints.meanCo2Value(0, 4)

	if mean != 11.5 {
		t.Error("unexpected mean", 11.5, mean)
	}
}

func TestReduce48Hours(t *testing.T) {
	startTime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

	rawPoints := make([]*models.SensorPoint, 0)

	for i := 0; i < 240; i++ {
		rawPoints = append(rawPoints, &models.SensorPoint{
			Time:     startTime.Add(-time.Minute * time.Duration(i)),
			Co2Value: 10.0 + float64(i),
		})
	}

	reducedPoints := NewReducedSensorPoints(rawPoints, startTime)

	fortyEightHours := reducedPoints.Last48Hours()

	if len(fortyEightHours) != 96 {
		t.Error("unexpected count", 96, len(fortyEightHours))
	}

	if !fortyEightHours[0].Time.Equal(startTime) {
		t.Error("unexpected first time value", startTime, fortyEightHours[0].Time)
	}

	if fortyEightHours[0].Co2Value != 24.5 {
		t.Error("unexpected first co2 value", 24.5, fortyEightHours[0].Co2Value)
	}

	if !fortyEightHours[7].Time.Equal(startTime.Add(-time.Minute * time.Duration(30*7))) {
		t.Error("unexpected 7th item time")
	}

	if fortyEightHours[7].Co2Value != 234.5 {
		t.Error("unexpected 7th item co2 value", 234.5, fortyEightHours[7].Co2Value)
	}

	if fortyEightHours[8].Co2Value != 400.0 {
		t.Error("unexpected 8th item co2 value", 400.0, fortyEightHours[8].Co2Value)
	}

	mean := reducedPoints.meanCo2Value(0, 4)

	if mean != 11.5 {
		t.Error("unexpected mean", 11.5, mean)
	}
}
