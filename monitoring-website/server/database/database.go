package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func ConnectDB() {
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

	db, err = sql.Open("mysql", cfg.FormatDSN())
	handleError(err)

	err = db.Ping()
	handleError(err)
}

func InsertEnergyData(dateTime string, rf int64) error {
	checkDBIsSet()

	_, err := db.Exec("INSERT INTO `energydata` (`dateTime`, `RF`) VALUES(?, ?)", dateTime, rf)
	return err
}

type DataPoint struct {
	DateTime string `json:"dateTime"`
	RF       int64  `json:"RF"`
}

func GetEnergyData(numEntries int64, startDateTime string) []DataPoint {
	checkDBIsSet()

	dateTime, err := db.Query("SELECT `dateTime` FROM `energydata` "+
		"ORDER BY ABS(TIMESTAMPDIFF(SECOND, `datetime`, ?)) "+
		"LIMIT 1",
		startDateTime)
	handleError(err)
	defer dateTime.Close()

	dateTime.Next()
	var closestDateTime string
	err = dateTime.Scan(&closestDateTime)
	handleError(err)

	var dataRows *sql.Rows
	dataRows, err = db.Query("SELECT `dateTime`, `RF` FROM `energydata` "+
		"WHERE `dateTime` <= ? "+
		"LIMIT ?",
		closestDateTime,
		numEntries)
	handleError(err)
	defer dataRows.Close()

	var dataSet []DataPoint
	for dataRows.Next() {
		dataPoint := DataPoint{}
		err = dataRows.Scan(&dataPoint.DateTime, &dataPoint.RF)
		dataSet = append(dataSet, dataPoint)
		handleError(err)
	}

	return dataSet
}

func checkDBIsSet() {
	if db == nil {
		ConnectDB()
	}
}

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
