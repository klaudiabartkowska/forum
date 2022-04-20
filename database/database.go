package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql  
)

func main() {

	// to open we need to provide (driverName:, dataSourceName:)

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	
	// 

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS people(id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	statement.Exec()

}
