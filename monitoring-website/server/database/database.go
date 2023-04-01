package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const MYSQL_TIME_FORMAT = "2006-01-02 15:04:05"

type Database struct {
	db *sql.DB
}

func ConnectDB() *Database {
	// Load env variables
	err := godotenv.Load()
	handleError(err)

	// MySQL configs
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DBADDR"),
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
	}

	database := &Database{}
	database.db, err = sql.Open("mysql", cfg.FormatDSN())
	handleError(err)

	err = database.db.Ping()
	handleError(err)

	return database
}

func (db *Database) Close() {
	db.db.Close()
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
