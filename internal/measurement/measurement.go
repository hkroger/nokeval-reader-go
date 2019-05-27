package measurement

import "time"

type Measurement struct {
	Timestamp      time.Time
	SensorId       int
	Measurement    float64
	Voltage        float64
	SignalStrength float64
}

func (m Measurement) Valid() bool {
	if m.Measurement < 200 && m.Voltage < 20 {
		return true
	}

	return false
}
