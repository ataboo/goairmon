package hardware

import "testing"

func TestFakeCo2SensorRead(t *testing.T) {
	sensor := NewFakeCo2Sensor(&Co2SensorCfg{}).(*fakeCo2Sensor)

	if err := sensor.On(); err != nil {
		t.Error(err)
	}

	for i := 0; i < 1000; i++ {
		val, _ := sensor.Read()
		if val < sensor.min || val > sensor.max {
			t.Error("value out of range", sensor.min, sensor.max, val)
		}
	}

}
