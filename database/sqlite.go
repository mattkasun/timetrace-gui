package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// SQLITE_FUNCTIONS - contains a map of the functions for sqlite
var SQLITE_FUNCTIONS = map[string]interface{}{
	INIT_DB:      sqInitDB,
	CREATE_TABLE: sqCreateTable,
	INSERT:       sqInsert,
	DELETE:       sqDeleteRecord,
	DELETE_ALL:   sqDeleteAllRecords,
	FETCH:        sqFetchRecords,
	CLOSE_DB:     sqCloseDB,
}

var db *sql.DB

func sqInitDB() error {
	//cfg, err := config.Get()
	//if err != nil {
	//log.Fatal("could not connect to database", err)
	//}
	//if cfg == (&config.Config{}) || cfg.DBFile == "" {
	//return errors.New("empty config file")
	//}
	//log.Println("initializing sqlite ", cfg.DBPath, cfg.DBFile)
	DBPath := "./"
	//DBPath := "/var/lib/timetrace/"
	DBFile := "timetrace.db"
	var err error

	if _, err := os.Stat(DBPath); os.IsNotExist(err) {
		if err := os.MkdirAll(DBPath, 0766); err != nil {
			log.Println("mkdir error: ", DBPath, err)
			return err
		}
	}
	path := filepath.Join(DBPath, DBFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			log.Println("file create error: ", err)
			return err
		}
	}
	db, err = sql.Open("sqlite", path)
	if err != nil {
		log.Println("error opening sqlite database", path, err)
		return err
	}
	db.SetMaxOpenConns(1)
	return db.Ping()
	//return nil
}

func sqCreateTable(table string) error {
	query := "CREATE TABLE IF NOT EXISTS " + table + " ( key TEXT NOT NULL UNIQUE PRIMARY KEY, value TEXT)"
	statement, err := db.Prepare(query)
	if err != nil {
		log.Println("error preparing query", err)
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Println("error executing statement", err)
		return err
	}
	return nil
}

func sqInsert(key, value, table string) error {
	//log.Println("sqlinsert", key, value, table)
	if key != "" && value != "" && json.Valid([]byte(value)) {
		insertSQL := "INSERT OR REPLACE INTO " + table + " (key, value) VALUES (?, ?)"
		statement, err := db.Prepare(insertSQL)
		if err != nil {
			return err
		}
		defer statement.Close()
		_, err = statement.Exec(key, value)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid insert " + key + " : " + value)
}

func sqDeleteRecord(id, table string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE KEY = '%s'", table, id)
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	if _, err := statement.Exec(); err != nil {
		return err
	}
	log.Printf("deleted %s from %s\n", id, table)
	return nil
}

func sqDeleteAllRecords() error {
	return nil
}

func sqFetchRecords(table string) (map[string]string, error) {
	query := "SELECT * FROM " + table + " ORDER BY key"
	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	records := make(map[string]string)
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var key string
		var value string
		row.Scan(&key, &value)
		records[key] = value
	}
	if len(records) == 0 {
		return nil, ErrNoResults
	}
	return records, nil
}

func sqCloseDB() {
	db.Close()
}
