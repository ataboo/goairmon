package hardware

import "testing"

func TestFakeCo2SensorRead(t *testing.T) {
	sensor := NewFakeSgp30Sensor().(*fakeSgp30)

	for i := 0; i < 1000; i++ {
		eCO2, TVOC, err := sensor.Measure()
		if err != nil {
			t.Error("unexpected error", err)
		}

		if eCO2 < sensor.min || eCO2 > sensor.max {
			t.Error("value out of range", sensor.min, sensor.max, eCO2)
		}

		if TVOC < sensor.min || TVOC > sensor.max {
			t.Error("value out of range", sensor.min, sensor.max, TVOC)
		}
	}
}
