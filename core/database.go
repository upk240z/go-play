package core

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Database struct {
	driver string
	dsn    string
	db     *sql.DB
	result sql.Result
}

func NewDatabase(driver, dsn string) *Database {
	instance := Database{driver, dsn, nil, nil}
	return &instance
}

func (database *Database) connect() {
	if database.db != nil {
		return
	}

	var err error
	database.db, err = sql.Open(database.driver, database.dsn)

	if err != nil {
		log.Fatal(err)
	}

	database.db.SetMaxOpenConns(0)
	database.db.SetMaxIdleConns(10)
	database.db.SetConnMaxLifetime(time.Minute * 5)
}

func (database *Database) All(query string, parameters ...interface{}) []map[string]*string {
	database.connect()
	rows, err := database.db.Query(query, parameters...)
	if err != nil {
		log.Fatal(err)
	}
	columns, _ := rows.Columns()

	var results []map[string]*string

	for rows.Next() {
		var pointers []interface{}
		mapValues := map[string]*string{}
		for _, columnName := range columns {
			var col string
			pointers = append(pointers, &col)
			mapValues[columnName] = &col
		}

		rows.Scan(pointers...)
		results = append(results, mapValues)
	}

	return results
}

func (database *Database) Row(query string, parameters ...interface{}) map[string]*string {
	database.connect()
	rows, err := database.db.Query(query, parameters...)
	if err != nil {
		log.Fatal(err)
	}
	columns, _ := rows.Columns()

	for rows.Next() {
		var pointers []interface{}
		mapValues := map[string]*string{}
		for _, columnName := range columns {
			var col string
			pointers = append(pointers, &col)
			mapValues[columnName] = &col
		}

		rows.Scan(pointers...)
		return mapValues
	}

	return nil
}

func (database *Database) Exec(query string, parameters ...interface{}) int64 {
	database.connect()

	var err1 error
	database.result, err1 = database.db.Exec(query, parameters...)
	if err1 != nil {
		log.Fatal(err1)
	}

	affected, err2 := database.result.RowsAffected()
	if err2 != nil {
		log.Fatal(err2)
	}

	return affected
}

func (database *Database) LastInsertId() int64 {
	if database.result == nil {
		return 0
	}

	lastId, err := database.result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return lastId
}

func (database *Database) Close() {
	if database.db == nil {
		return
	}

	if err := database.db.Close(); err != nil {
		log.Fatal(err)
	}
}
