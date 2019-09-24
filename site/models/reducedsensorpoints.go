package models

import (
	"goairmon/business/data/models"
	"sort"
	"time"
)

func NewReducedSensorPoints(rawPointData []*models.SensorPoint, now time.Time) *ReducedSensorPoints {
	reduced := &ReducedSensorPoints{}
	reduced.normalizeSensorData(rawPointData, now)

	return reduced
}

type ReducedSensorPoints struct {
	pointData []*models.SensorPoint
}

func (p *ReducedSensorPoints) Last7Days() []*models.SensorPoint {
	// Mean point by 60 minutes
	outputPoints := 7 * 24
	pointRange := 60
	output := make([]*models.SensorPoint, outputPoints)

	for i := 0; i < outputPoints; i++ {
		midPointIdx := pointRange * i
		midPointTime := p.pointData[midPointIdx].Time
		meanCo2 := p.meanCo2Value(midPointIdx, pointRange)

		output[i] = &models.SensorPoint{
			Time:     midPointTime,
			Co2Value: meanCo2,
		}
	}

	return output
}

func (p *ReducedSensorPoints) Last48Hours() []*models.SensorPoint {
	// Mean point by 30 minutes
	outputPoints := 24 * 2 * 2
	pointRange := 30
	output := make([]*models.SensorPoint, outputPoints)

	for i := 0; i < outputPoints; i++ {
		midPointIdx := pointRange * i
		midPointTime := p.pointData[midPointIdx].Time
		meanCo2 := p.meanCo2Value(midPointIdx, pointRange)

		output[i] = &models.SensorPoint{
			Time:     midPointTime,
			Co2Value: meanCo2,
		}
	}

	return output
}

func (p *ReducedSensorPoints) Last2Hours() []*models.SensorPoint {
	pointCount := 120
	output := make([]*models.SensorPoint, pointCount)

	for i := 0; i < pointCount; i++ {
		point := p.pointData[i]

		output[i] = &models.SensorPoint{
			Time:     point.Time,
			Co2Value: point.Co2Value,
		}
	}

	return output
}

func (p *ReducedSensorPoints) meanCo2Value(minPointIdx int, pointRange int) float64 {
	sum := 0.0
	for i := minPointIdx; i < minPointIdx+pointRange; i++ {
		sum += p.pointData[i].Co2Value
	}

	return sum / float64(pointRange)
}

func (p *ReducedSensorPoints) normalizeSensorData(rawPoints []*models.SensorPoint, now time.Time) {
	pointCount := 24 * 8 * 60
	p.pointData = make([]*models.SensorPoint, pointCount)

	sort.Slice(rawPoints, func(i, j int) bool {
		return rawPoints[i].Time.After(rawPoints[j].Time)
	})

	rawIdx := 0
	var refTime time.Time
	nextRawPoint := rawPoints[rawIdx]

	for i := 0; i < pointCount; i++ {
		refTime = now.Add(-time.Minute * time.Duration(i))

		for nextRawPoint != nil && nextRawPoint.Time.After(refTime) {
			rawIdx++
			if rawIdx >= len(rawPoints) {
				nextRawPoint = nil
			} else {
				nextRawPoint = rawPoints[rawIdx]
			}
		}

		co2Value := 400.0
		if nextRawPoint != nil && nextRawPoint.Time.Add(time.Minute).After(refTime) {
			co2Value = nextRawPoint.Co2Value
		}

		p.pointData[i] = &models.SensorPoint{
			Time:     refTime,
			Co2Value: co2Value,
		}
	}
}
