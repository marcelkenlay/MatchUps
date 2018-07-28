package handlers

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"

	. "../utils"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
	DB_HOST     = "35.176.172.232"
	DB_PORT     = "5432"
)

var Database *sql.DB

func ConnectToDatabase() {
	// Set up connection
	if (Database != nil) {
		return
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
	db, err := sql.Open("postgres", dbinfo)
	Database = db
	CheckErr(err)
}
