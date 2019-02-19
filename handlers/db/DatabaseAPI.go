package db

import (
	. "../utils"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"os"
)

var Database *sql.DB


func ConnectToDatabase() {
	// Set up connection
	if Database != nil {
		return
	}

	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")

	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DbUser, DbPassword, DbName, DbHost, DbPort)

	print(dbInfo)

	db, err := sql.Open("postgres", dbInfo)
	Database = db
	CheckErr(err)
}

func QueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}