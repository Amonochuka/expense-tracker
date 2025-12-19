package db

import (
	"database/sql"
	"log"

	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string, sqlFilePath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB : %v", err)
	}

	//ensure connection starts
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping DB : %v", err)
	}

	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("Failed to read SQL file :%v", err)
	}

	sqlStatement := string(sqlBytes)
	_, err = DB.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Failed to execute SQL file: %v", err)
	}

	log.Println("Database initialized successfully")
}
