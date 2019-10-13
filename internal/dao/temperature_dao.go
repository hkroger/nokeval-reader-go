package dao

import (
	"bytes"
	"crypto"
	"encoding/json"
	"fmt"
	"github.com/hkroger/nokeval-reader-go/internal/config"
	"github.com/hkroger/nokeval-reader-go/internal/measurement"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type MeasurementDAO struct {
	DatabaseConfig config.DatabaseConfig
}

func ftos(n float64) string {
	return fmt.Sprintf("%f", n)
}

func (dao MeasurementDAO) Store(reading measurement.Measurement) error {
	content := make(map[string]interface{})

	content["client_id"] = dao.DatabaseConfig.ClientId
	content["timestamp"] = fmt.Sprintf("%d", reading.Timestamp.Unix())
	content["sensor_id"] = strconv.Itoa(reading.SensorId)
	content["measurement"] = ftos(reading.Measurement)
	content["voltage"] = ftos(reading.Voltage)
	content["signal_strength"] = ftos(reading.SignalStrength)
	content["version"] = 2
	content["checksum"] = generateChecksum(content, dao.DatabaseConfig.Secret)

	for _, url := range dao.DatabaseConfig.OverrideUrls {
		jsonValue, err := json.Marshal(content)

		if err != nil {
			log.Fatalf("Could not marshal json: %s", err)
		}

		request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			return fmt.Errorf("Technical error when trying to speak HTTP: %s", err)
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			log.Debugf("Measurement stored in %s successfully", url)
			return nil
		} else if response.StatusCode == 403 {
			log.Debug("Got 403, we are not authorized. Let's skip this.")
			return nil
		}

		body, _ := ioutil.ReadAll(response.Body)

		log.Errorf("Could not store the result. Error Code: %s, Response: %s", response.StatusCode, body)
	}

	return nil
}

func generateChecksum(hsh map[string]interface{}, secret string) string {
	src := fmt.Sprintf("%d&%s&%s&%s&%s&%s&%s&%s",
		hsh["version"],
		hsh["timestamp"],
		hsh["voltage"],
		hsh["signal_strength"],
		hsh["client_id"],
		hsh["sensor_id"],
		hsh["measurement"],
		secret,
	)

	log.Debugf("Checksum source string: %s", src)

	hsher := crypto.SHA1.New()
	_, _ = io.WriteString(hsher, src)
	return fmt.Sprintf("%x", hsher.Sum(nil))
}
