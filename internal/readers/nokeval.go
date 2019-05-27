package readers

import (
	"encoding/hex"
	"errors"
	"github.com/hkroger/nokeval-reader-go/internal/config"
	"github.com/hkroger/nokeval-reader-go/internal/measurement"
	"github.com/jacobsa/go-serial/serial"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
	time2 "time"
)

var address uint8 = 0
var noValues uint8 = '#'

var deviceTypes = map[int]string{
	0:  "MTR260",
	2:  "MTR262",
	4:  "MTR264",
	5:  "MTR265",
	6:  "MTR165",
	7:  "FTR860",
	8:  "CSR264S",
	9:  "CSR264L",
	10: "CSR264A",
	11: "CSR260",
	12: "KMR260",
}

var id uint8 = 128
var etx uint8 = 3
var ack uint8 = 6

type NokevalTemperatureReader struct {
	SerialConfig config.SerialConfig
	port         io.ReadWriteCloser
}

func (r *NokevalTemperatureReader) Next() (*measurement.Measurement, error) {
	r.sclCommand("DBG 1 ?", address)
	response := r.sclResponse()

	log.Debugf("Response: %s", response)

	if response[0] != noValues {
		split := strings.Split(response, " ")
		if split[0] == "0" {
			devType, error := strconv.Atoi(split[0])
			if error != nil {
				return nil, error
			}

			rawVoltage, error := strconv.Atoi(split[1])
			if error != nil {
				return nil, error
			}
			voltage := float64(rawVoltage&31) / 10.0

			rawSignalStrength, error := strconv.Atoi(split[2])
			if error != nil {
				return nil, error
			}

			signalStrength := float64(rawSignalStrength&127) - 127

			deviceId, error := strconv.Atoi(split[3])
			if error != nil {
				return nil, error
			}

			measurementPart1, error := strconv.Atoi(split[4])
			if error != nil {
				return nil, error
			}
			measurementPart2, error := strconv.Atoi(split[5])
			if error != nil {
				return nil, error
			}
			measurementValue := float64(measurementPart1+measurementPart2*256)/10.0 - 273.2
			time := time2.Now()
			log.Debugf("%s: Device %d, Device Type %s, Signal strength %fdBm, voltage: %fV, measurement: %f",
				time,
				deviceId,
				deviceTypes[devType],
				signalStrength,
				voltage,
				measurementValue,
			)
			return &measurement.Measurement{Timestamp: time, SensorId: deviceId, Measurement: measurementValue, Voltage: voltage, SignalStrength: signalStrength}, nil
		}
	}
	return nil, errors.New("No worving")
}

func calcBcc(str []byte) uint8 {
	var checksum uint8 = 0
	for _, char := range str {
		checksum = checksum ^ uint8(char)
	}
	return checksum
}

func (r *NokevalTemperatureReader) sclCommand(cmd string, address uint8) {
	var message = make([]byte, len(cmd)+3)

	message[0] = address + id
	for i, char := range cmd {
		message[i+1] = uint8(char)
	}
	message[len(cmd)+1] = etx
	message[len(cmd)+2] = calcBcc(message[1:])

	log.Tracef("Writing SCL command: %s", cmd)
	log.Tracef("message: %s", message)
	log.Tracef("Length: %d", len(message))
	log.Tracef("Hexdump: %s", hex.Dump([]byte(message)))
	n, error := r.port.Write([]byte(message))

	if error != nil {
		log.Errorf("Writing SCL command failed: %s. Wrote %d bytes out of %d", error, n, len(message))
	} else {
		log.Tracef("Writing SCL command successful.")

	}
}

func (r *NokevalTemperatureReader) sclResponse() string {
	buffer := ""
	for {
		shortReadBuffer := make([]byte, 1)

		log.Trace("Reading 1 byte from serial port")
		n, err := r.port.Read(shortReadBuffer)
		log.Trace("Reading 1 byte from serial port done")

		if err == nil && n == 1 {
			buffer = buffer + string(shortReadBuffer)
			msg := sclValidResponse(buffer)

			if msg != nil {
				return *msg
			}
		} else if err != nil {
			log.Debugf("Error reading serial port: %s", err)
		}
	}
}

func sclValidResponse(message string) *string {
	// size := len(message)
	if message[0] != ack {
		return nil
	}

	etxOffset := len(message) - 2

	if etxOffset <= 0 {
		return nil
	}

	if message[etxOffset] != etx {
		return nil
	}

	bccOffset := len(message) - 1
	calculatedBcc := calcBcc([]byte(message[0 : len(message)-1]))

	if message[bccOffset] != calculatedBcc {
		return nil
	}

	var validMsg = message[1 : len(message)-2]

	return &validMsg
}

func (r *NokevalTemperatureReader) Open() error {
	port, err := serial.Open(r.openConfig())

	if err != nil {
		return err
	}

	r.port = port

	return nil
}

func (r NokevalTemperatureReader) openConfig() serial.OpenOptions {
	return serial.OpenOptions{
		PortName:        r.SerialConfig.Device,
		BaudRate:        r.SerialConfig.Baud,
		DataBits:        r.SerialConfig.Bits,
		StopBits:        r.SerialConfig.StopBits,
		MinimumReadSize: 1,
	}

}

func (r *NokevalTemperatureReader) close() {
	if r.port != nil {
		err := r.port.Close()

		if err != nil {
			log.Errorf("Closing port failed: %s", err)
		}
		r.port = nil
	}
}
