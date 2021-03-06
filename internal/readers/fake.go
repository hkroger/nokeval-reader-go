package readers

import (
	"github.com/hkroger/nokeval-reader-go/internal/measurement"
	"math/rand"
	"time"
)

type FakeTemperatureReader struct {
}

func (FakeTemperatureReader) Next() (*measurement.Measurement, error) {
	time.Sleep(1 * time.Millisecond)
	if rand.Intn(1000) > 0 {
		return &measurement.Measurement{
			Timestamp:      time.Now(),
			SensorId:       600000 + rand.Intn(8),
			Measurement:    rand.NormFloat64() * 30,
			Voltage:        3.3 + rand.NormFloat64()*1,
			SignalStrength: rand.NormFloat64() * 100,
		}, nil
	} else {
		return nil, nil
	}
}

func (FakeTemperatureReader) Open() error {
	return nil
}
