package core

import (
	"database/sql"
	"log"
	"regexp"
	"strings"
	"time"
)

type Database struct {
	driver              string
	dsn                 string
	db                  *sql.DB
	result              sql.Result
	query               string
	parameters          map[string]interface{}
	convertedParameters []interface{}
}

func NewDatabase(driver, dsn string) *Database {
	instance := Database{driver: driver, dsn: dsn}
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

func (database *Database) convertNamedPlaceHolder() {
	pattern, _ := regexp.Compile(`:([a-z\d\-_]+)`)

	database.convertedParameters = []interface{}{}

	for {
		if !pattern.MatchString(database.query) {
			break
		}

		mark := pattern.FindString(database.query)
		p, exists := database.parameters[mark[1:]]
		if !exists {
			log.Fatal("parameter not found: " + mark[1:])
		}

		database.convertedParameters = append(database.convertedParameters, p)
		database.query = strings.Replace(database.query, mark, "?", 1)
	}
}

func (database *Database) queryCommon(query string, parameters map[string]interface{}) *sql.Rows {
	database.connect()
	database.query = query
	database.parameters = parameters
	database.convertNamedPlaceHolder()

	rows, err := database.db.Query(database.query, database.convertedParameters...)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func (database *Database) All(query string, parameters map[string]interface{}) []map[string]*string {
	rows := database.queryCommon(query, parameters)
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

		err := rows.Scan(pointers...)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, mapValues)
	}

	return results
}

func (database *Database) Row(query string, parameters map[string]interface{}) map[string]*string {
	rows := database.queryCommon(query, parameters)
	columns, _ := rows.Columns()

	for rows.Next() {
		var pointers []interface{}
		mapValues := map[string]*string{}
		for _, columnName := range columns {
			var col string
			pointers = append(pointers, &col)
			mapValues[columnName] = &col
		}

		err := rows.Scan(pointers...)
		if err != nil {
			log.Fatal(err)
		}

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
