package readers

import "github.com/hkroger/nokeval-reader-go/internal/measurement"

type TemperatureReader interface {
	Next() (*measurement.Measurement, error)
	Open() error
}
