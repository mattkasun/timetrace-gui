package database

import (
	"encoding/json"
	"errors"

	"github.com/mattkasun/timetrace-gui/models"
)

func Saverecord(record *models.Record) error {
	value, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return insert(record.ID.String(), string(value), RECORDS_TABLE_NAME)
}

func GetRecord(id string) (models.Record, error) {
	var record models.Record
	records, err := fetch(RECORDS_TABLE_NAME)
	if err != nil {
		return record, err
	}
	for key, value := range records {
		if key == id {
			if err := json.Unmarshal([]byte(value), &record); err != nil {
				return record, err
			}
			return record, nil
		}
	}
	return record, errors.New("no such record")
}

func GetAllrecords() ([]models.Record, error) {
	var records []models.Record
	var record models.Record
	rows, err := fetch(RECORDS_TABLE_NAME)
	if err != nil {
		return records, err
	}
	for _, value := range rows {
		if err := json.Unmarshal([]byte(value), &record); err != nil {
			continue
		}
		records = append(records, record)
	}
	return records, nil
}

func DeleteRecord(id string) error {
	return delete(id, RECORDS_TABLE_NAME)
}
