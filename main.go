package main

import (
	"database/sql"
	"fmt"

	"net/http"

	"example.com/m/database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./database/userdata.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	data := database.Connect(db)

	// data.CreateUser(database.User{})

	//defer DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", data.Handler)   // this is not a package it's a method
	fmt.Println("Starting the server on :8080...")
	if err := http.ListenAndServe(":8080", mux)
	 err != nil {
		panic(err)
	}

}
