package buffer

import (
	"database/sql"
	json2 "encoding/json"
	"fmt"
	"github.com/hkroger/nokeval-reader-go/internal/dao"
	"github.com/hkroger/nokeval-reader-go/internal/measurement"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type MeasurementBuffer struct {
	BufferFile string
	connection *sql.DB
}

func (buffer *MeasurementBuffer) Flush(dao *dao.MeasurementDAO) error {
	for {
		var rowsProcessed, err = buffer.doFlush(dao)
		if err != nil {
			return err
		}

		if !rowsProcessed {
			break
		}
	}

	return nil
}

func daoStore(dao *dao.MeasurementDAO, requestChannel chan request, responseChannel chan response, processorFinished chan bool) {
	for {
		request, ok := <- requestChannel

		// Channel closed
		if !ok {
			log.Debugf("Channel closed")
			break
		}

		err := dao.Store(request.measurement)

		if err != nil {
			log.Errorf("DAO storage failed: %s", err)
		}

		responseChannel <- response{id:request.id, error:err}
	}

	processorFinished <- true
}

type request struct {
	id int
	measurement measurement.Measurement
}

type response struct {
	id int
	measurement measurement.Measurement
	error error
}

func (buffer *MeasurementBuffer) doFlush(dao *dao.MeasurementDAO) (bool, error) {
	var channelSize = 10
	var processors = 5
	rows, err := buffer.connection.Query("SELECT id, json_data FROM measurements LIMIT ?", channelSize)
	var doneCounter = 0

	storeRequests := make(chan request, channelSize)
	storeResponses := make(chan response, channelSize)
	processorFinished := make(chan bool, processors)

	for i:=0; i<processors; i++ {
		go daoStore(dao, storeRequests, storeResponses, processorFinished)
	}

	if err != nil {
		return false, fmt.Errorf("Query failed: %v, connection: %v", err, buffer.connection)
	}
	defer rows.Close()

	for rows.Next() {
		// Use Scan to access column data from a row
		var id int
		var json string
		var measurement measurement.Measurement
		err = rows.Scan(&id, &json)

		if err == nil {
			err = json2.Unmarshal([]byte(json), &measurement)

			if err != nil {
				return false, fmt.Errorf("Could not deserialize json from buffer: %s, json: %s", err, json)
			}

			storeRequests <- request{id:id, measurement:measurement}
		} else {
			log.Errorf("scan failed while querying measurements: %v", err)
			err = nil // Let's ignore this and let the problematic line to be deleted.
		}
	}

	close(storeRequests)

	for {
		<- processorFinished
		doneCounter++

		if doneCounter >= processors {
			close(processorFinished)
			break
		}
	}

	close(storeResponses)

	deletedRows := 0

	for {
		response, ok := <- storeResponses

		if !ok {
			break
		}

		err = buffer.remove(response.id)
		if err != nil {
			log.Fatalf("Deletion of measurement from buffer failed: %s", err)
		}
		deletedRows++
	}

	return deletedRows>0, nil
}

func (buffer *MeasurementBuffer) Store(measurement *measurement.Measurement) error {
	data, err := json2.Marshal(measurement)
	if err != nil {
		return err
	}

	_, err = buffer.connection.Exec("INSERT INTO measurements(json_data) VALUES (?)", data)
	return err
}

func (buffer *MeasurementBuffer) remove(id int) error {
	_, err := buffer.connection.Exec("DELETE FROM measurements WHERE id = ?", id)
	return err
}

func (buffer *MeasurementBuffer) Open() error {
	if buffer.BufferFile == "" {
		log.Fatalf("buffer_file configuration is empty :(")
	}

	log.Debugf("Opening DB file: %s", buffer.BufferFile)
	conn, err := sql.Open("sqlite3", buffer.BufferFile)
	if err != nil {
		return err
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS measurements( id INTEGER PRIMARY KEY ASC, json_data TEXT)")

	if err != nil {
		return err
	}

	buffer.connection = conn

	log.Debugf("sqlite connection: %v", buffer.connection)

	return nil
}
