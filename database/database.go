package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" //connection driver
)

// Connect to db
func Connect() (*sql.DB, error) {
	connectionString := "<username>:<password>@/<database_name>?charset=utf8&parseTime=True&loc=Local"

	db, error := sql.Open("mysql", connectionString)
	if error != nil {
		log.Fatal(error)
		return nil, error
	}
	if error = db.Ping(); error != nil {

		log.Fatal(error)
		return nil, error
	}

	return db, nil
}
