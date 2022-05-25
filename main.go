package main

import (
	"database/sql"
	"fmt"

	"net/http"

	"example.com/m/database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	data := database.Connect(db)

	// data.CreateUser(database.User{})

	// data.CreateUser(database.User{
	// 	Username: "Klaudia",
	// 	Email:    "sasa",
	// 	Password: "Barti",
	// })

	//defer DB.Close()

	//   database.Connect()



	mux := http.NewServeMux()
	mux.HandleFunc("/", data.Handler)   // this is not a package it's a method

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	fmt.Println("Starting the server on :8080...")
	server.ListenAndServe()

	// fmt.Println("Starting the server on :8080...")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	panic(err)
	// }

}
