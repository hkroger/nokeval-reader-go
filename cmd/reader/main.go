package main

import (
	"flag"
	"github.com/hkroger/nokeval-reader-go/internal/buffer"
	"github.com/hkroger/nokeval-reader-go/internal/config"
	"github.com/hkroger/nokeval-reader-go/internal/dao"
	"github.com/hkroger/nokeval-reader-go/internal/measurement"
	"github.com/hkroger/nokeval-reader-go/internal/readers"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

func main() {
	log.Info("Starting nokeval reader...")

	appConfig := initConfig()

	measurementDao := dao.MeasurementDAO{MeasurementStorageConfig: appConfig.MeasurementStorage}
	measurementBuffer := buffer.MeasurementBuffer{BufferFile: appConfig.BufferFile}
	err := measurementBuffer.Open()

	if err != nil {
		log.Fatalf("Could not open measurement buffer: %s", err)
	}

	flushMeasurements(&measurementBuffer, &measurementDao, &appConfig)

	for {
		var reader readers.TemperatureReader
		if appConfig.FakeSensorMode {
			reader = &readers.FakeTemperatureReader{}
		} else {
			reader = &readers.NokevalTemperatureReader{SerialConfig: appConfig.Serial}
		}

		err := reader.Open()

		if err != nil {
			log.Errorf("Reader open failed: %v", err)
			time.Sleep(10 * time.Second)
		} else {
			for {
				log.Debug("Reading next measurement")
				reading, err := reader.Next()
				if err == nil && reading != nil && reading.Valid() {
					storeMeasurement(&appConfig, reading, &measurementBuffer)
				}

				if reading == nil {
					break
				}
			}

			flushMeasurements(&measurementBuffer, &measurementDao, &appConfig)

			log.Debug("Waiting 5 seconds")
			time.Sleep(5 * time.Second)
		}
	}
}

func storeMeasurement(appConfig *config.Config, reading *measurement.Measurement, measurementBuffer *buffer.MeasurementBuffer) {
	if appConfig.FakeStorageMode {
		log.Debugf("fake storage mode: %v", reading)
	} else {
		log.Debugf("storing measurement: %v", reading)

		err := measurementBuffer.Store(reading)

		if err != nil {
			log.Panicf("Buffer storage failed. Let's bail out. Error: %v", err)
		}
	}
}

func flushMeasurements(measurementBuffer *buffer.MeasurementBuffer, measurementDao *dao.MeasurementDAO, appConfig *config.Config) {
	if !appConfig.FakeStorageMode {

		log.Debug("Flushing measurements")
		err := measurementBuffer.Flush(measurementDao)
		if err != nil {
			log.Errorf("Could not flush: %s", err)
		}
	}
}

func initConfig() config.Config {
	configPath := flag.String("c", "/etc/nokeval-reader.yaml", "Config file path, defaults to /etc/nokeval-reader.yaml")
	debug := flag.Bool("v", false, "Verbose mode, defaults to false")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Config file: %s", *configPath)
	configYaml, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Panicf("Failed to read config file: %s, error: %s", *configPath, err)
	}

	c := config.Config{}
	err = yaml.Unmarshal(configYaml, &c)
	if err != nil {
		log.Panicf("Failed to parse config file contents: %s, error: %s", configYaml, err)
	}
	log.Debugf("Config file contents: %v", c)

	if len(c.MeasurementStorage.OverrideUrls) <= 0 {
		c.MeasurementStorage.OverrideUrls = []string{"https://api.measurinator.com/measurements"}
	}

	return c
}
