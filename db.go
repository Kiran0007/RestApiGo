package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDb() {

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_SERVER := os.Getenv("DB_SERVER")
	DB_PORT := "3306"
	DB_NAME := os.Getenv("DB_NAME")

	if DB_USER == "" {
		log.Fatal("DB_USER is not set as Environment variable")
	} else if DB_PASSWORD == "" {
		log.Fatal("DB_PASSWORD is not set as Environment variable")
	} else if DB_SERVER == "" {
		log.Fatal("DB_SERVER is not set as Environment variable")
	} else if DB_NAME == "" {
		log.Fatal("DB_NAME is not set as Environment variable")
	}

	var err error
	db, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp("+DB_SERVER+":"+DB_PORT+")/"+DB_NAME)
	if err != nil {
		log.Fatal("cannot initialize db", err)
	} else {
		log.Println("SQL connection opened")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("cannot ping db!")
	} else {
		log.Println("Success pinging db.")
	}
}
