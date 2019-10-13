package buffer

import (
	json2 "encoding/json"
	"fmt"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/hkroger/nokeval-reader-go/internal/dao"
	"github.com/hkroger/nokeval-reader-go/internal/measurement"
	log "github.com/sirupsen/logrus"
)

type MeasurementBuffer struct {
	BufferFile string
	connection *sqlite3.Conn
}

func (buffer *MeasurementBuffer) Flush(dao *dao.MeasurementDAO) error {
	stmt, err := buffer.connection.Prepare("SELECT id, json_data FROM measurements")
	if err != nil {
		return fmt.Errorf("prepare failed: %v, connection: %v", err, buffer.connection)
	}
	defer stmt.Close()

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return fmt.Errorf("step failed while querying students: %v", err)
		}
		if !hasRow {
			break
		}

		// Use Scan to access column data from a row
		var id int
		var json string
		var measurement measurement.Measurement
		err = stmt.Scan(&id, &json)

		if err == nil {
			err = json2.Unmarshal([]byte(json), &measurement)

			if err != nil {
				return fmt.Errorf("scan failed while querying measurements: %v", err)
			}

			err = dao.Store(measurement)
		} else {
			log.Errorf("Could not deserialize json from buffer: %s, json: %s", err, json)
			err = nil // Let's ignore this and let the problematic line to be deleted.
		}

		if err == nil {
			err = buffer.remove(id)
			if err != nil {
				log.Fatalf("Deletion of measurement from buffer failed: %s", err)
			}
		} else {
			log.Errorf("DAO storage failed: %s", err)
		}
	}

	return nil
}

func (buffer *MeasurementBuffer) Store(measurement *measurement.Measurement) error {
	data, err := json2.Marshal(measurement)
	if err != nil {
		return err
	}

	err = buffer.connection.Exec("INSERT INTO measurements(json_data) VALUES (?)", data)
	return err
}

func (buffer *MeasurementBuffer) remove(id int) error {
	return buffer.connection.Exec("DELETE FROM measurements WHERE id = ?", id)
}

func (buffer *MeasurementBuffer) Open() error {
	conn, err := sqlite3.Open(buffer.BufferFile)
	if err != nil {
		return err

	}

	err = conn.Exec("CREATE TABLE IF NOT EXISTS measurements( id INTEGER PRIMARY KEY ASC, json_data TEXT )")

	if err != nil {
		return err
	}

	buffer.connection = conn

	log.Debugf("sqlite connection: %v", buffer.connection)

	return nil
}
