package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

var DB *sql.DB

func Connect() {

	// to open we need to provide (driverName:, dataSourceName:)

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		fmt.Println(err.Error())
	}

	//defer db.Close()

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS people(id INTEGER PRIMARY KEY, username TEXT, email TEXT, passwordHash TEXT)")
	statement.Exec()

	DB = db
}
