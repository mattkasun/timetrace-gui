package database

import (
	"log"
)

const (
	// table names
	USERS_TABLE_NAME   = "users"
	PROJECT_TABLE_NAME = "projects"
	RECORDS_TABLE_NAME = "records"
	// sql verbs
	INIT_DB      = "init"
	CREATE_TABLE = "createtable"
	INSERT       = "insert"
	DELETE       = "delete"
	DELETE_ALL   = "deleteall"
	FETCH        = "fetch"
	CLOSE_DB     = "close"
	// errors
	NO_RECORDS = "no results found"
)

func getCurrentDB() map[string]interface{} {
	//config, _ := config.Get()
	//switch config.DB {
	//case "sqlite":
	return SQLITE_FUNCTIONS
	//case "postgres":
	//	return POSTGRES_FUNCTIONS
	//default:
	//	return SQLITE_FUNCTIONS
	//}
}

func InitializeDatabase() error {
	log.Println("connecting to database")
	if err := getCurrentDB()[INIT_DB].(func() error)(); err != nil {
		return err
	}
	return createTables()
}

func createTables() error {
	if err := createTable(USERS_TABLE_NAME); err != nil {
		return err
	}
	if err := createTable(PROJECT_TABLE_NAME); err != nil {
		return err
	}
	if err := createTable(RECORDS_TABLE_NAME); err != nil {
		return err
	}
	return nil
}

func createTable(name string) error {
	return getCurrentDB()[CREATE_TABLE].(func(string) error)(name)
}

func insert(key, value, table string) error {
	return getCurrentDB()[INSERT].(func(string, string, string) error)(key, value, table)
}

func fetch(table string) (map[string]string, error) {
	return getCurrentDB()[FETCH].(func(string) (map[string]string, error))(table)
}

func delete(key, table string) error {
	return getCurrentDB()[DELETE].(func(string, string) error)(key, table)
}
