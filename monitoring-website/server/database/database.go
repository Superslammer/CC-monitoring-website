package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const MYSQL_TIME_FORMAT = "2006-01-02 15:04:05"

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

func InsertEnergyData(dateTime string, rf int64, computerID int) error {
	checkDBIsSet()

	_, err := db.Exec("INSERT INTO `energydata` (`dateTime`, `RF`, `computerID`) VALUES(?, ?, ?)", dateTime, rf, computerID)
	if err != nil {
		return err
	}

	// Update `lastUpdated` in energyComputers
	curTime := time.Now().Format(MYSQL_TIME_FORMAT)
	_, err = db.Exec("UPDATE `computers` SET `lastUpdate`=? WHERE `id`=?", curTime, computerID)
	return err
}

type DataPoint struct {
	DateTime   string `json:"dateTime"`
	RF         int64  `json:"RF"`
	ComputerID int    `json:"computerID"`
}

func GetEnergyData(computerID int, numEntries int, startDateTime string) ([]DataPoint, error) {
	checkDBIsSet()

	dateTime := db.QueryRow("SELECT `dateTime` FROM `energydata` "+
		"ORDER BY ABS(TIMESTAMPDIFF(SECOND, `datetime`, ?)) "+
		"LIMIT 1", startDateTime)

	err := dateTime.Err()
	if err != nil {
		nilData := []DataPoint{}
		err = errors.New("No entries found")
		return nilData, err
	}
	var closestDateTime string
	err = dateTime.Scan(&closestDateTime)
	handleError(err)

	var dataRows *sql.Rows
	if computerID == -1 {
		dataRows, err = db.Query("SELECT `dateTime`, `RF`, `computerID` "+
			"FROM `energydata` "+
			"WHERE `dateTime` <= ? "+
			"ORDER BY `dateTime` DESC "+
			"LIMIT ?",
			closestDateTime,
			numEntries)
		handleError(err)
	} else {
		dataRows, err = db.Query("SELECT `dateTime`, `RF`, `computerID` "+
			"FROM `energydata` "+
			"WHERE `computerID` = ? AND `dateTime` <= ? "+
			"ORDER BY `dateTime` DESC "+
			"LIMIT ?",
			computerID, closestDateTime, numEntries)
		handleError(err)
	}
	defer dataRows.Close()

	var dataSet []DataPoint
	for dataRows.Next() {
		var dataPoint DataPoint
		err = dataRows.Scan(&dataPoint.DateTime, &dataPoint.RF, &dataPoint.ComputerID)
		handleError(err)
		dataSet = append(dataSet, dataPoint)
	}
	return dataSet, nil
}

func CreateEnergyComputerEntry(id int, maxEnergy int64, name string) error {
	checkDBIsSet()

	// Check if computer is already registered
	energyComputerRow, err := db.Query("SELECT `id` FROM `energyComputers` WHERE `id`=? LIMIT 1", id)
	handleError(err)
	defer energyComputerRow.Close()
	if energyComputerRow.Next() {
		return errors.New("Computer already registered as a energy computer")
	}

	var computerRow *sql.Rows
	computerRow, err = db.Query("SELECT `typeID` FROM `computers` WHERE `id`=? LIMIT 1", id)
	handleError(err)
	defer computerRow.Close()

	var typeID int
	if computerRow.Next() {
		computerRow.Scan(&typeID)
		curTime := time.Now().Format(MYSQL_TIME_FORMAT)
		_, err = db.Exec("UPDATE `computers` SET `typeName`=?, `name`=?, `lastUpdate`=? WHERE `id`=?", 1, name, curTime, id)
		handleError(err)
	} else {
		RegisterComputer(id, 1, name)
	}

	_, err = db.Exec("INSERT INTO `energycomputers` (`id`, `maxRF`) VALUES(?, ?)", id, maxEnergy)
	handleError(err)

	return nil
}

func RegisterComputer(id int, typeID int, name string) {
	checkDBIsSet()

	curTime := time.Now().Format(MYSQL_TIME_FORMAT)
	_, err := db.Exec("INSERT INTO `computers` (`id`, `typeID`, `name`, `lastUpdate`) VALUES(?, ?, ?, ?)", id, typeID, name, curTime)
	handleError(err)
}

type EnergyComputerResponse struct {
	Error        bool                            `json:"error"`
	ComputerData []EnergyComputerResponseElement `json:"computers"`
}

type EnergyComputerResponseElement struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MaxRF       int64  `json:"maxRF"`
	LastUpdated string `json:"lastUpdated"`
	CurrentRF   int64  `json:"currentRF"`
}

// TODO: Check if any computers have been regisered and return error if so
func GetEnergyComputers(id, numComputers int) (EnergyComputerResponse, error) {
	var computers EnergyComputerResponse

	if id == -1 {
		energyComputerRows, err := db.Query("SELECT `id`, `maxRF` FROM `energycomputers` LIMIT ?", numComputers)
		handleError(err)
		defer energyComputerRows.Close()

		for energyComputerRows.Next() {
			var computerElement EnergyComputerResponseElement
			var computerID int
			var maxRF int64
			err = energyComputerRows.Scan(&computerID, &maxRF)
			handleError(err)

			var name string
			var lastUpdate string
			computerRow := db.QueryRow("SELECT `name`, `lastUpdate` FROM `computers` WHERE `id`=? LIMIT 1", computerID)
			err = computerRow.Scan(&name, &lastUpdate)
			handleError(err)

			var currentRF int64
			energyDataRow := db.QueryRow("SELECT `RF` FROM `energydata` WHERE `computerID`=? ORDER BY `dateTime` DESC LIMIT 1", computerID)
			err = energyDataRow.Scan(&currentRF)
			handleError(err)

			computerElement = EnergyComputerResponseElement{
				ID:          computerID,
				Name:        name,
				MaxRF:       maxRF,
				LastUpdated: lastUpdate,
				CurrentRF:   currentRF,
			}
			computers.ComputerData = append(computers.ComputerData, computerElement)
		}
	} else {
		var computerElement EnergyComputerResponseElement
		var maxRF int64
		energyComputerRow := db.QueryRow("SELECT `maxRF` FROM `energycomputers` WHERE `id`=? LIMIT 1", id)
		err := energyComputerRow.Scan(&maxRF)
		handleError(err)

		var name string
		var lastUpdate string
		computerRow := db.QueryRow("SELECT `name`, `lastUpdate` FROM `computers` WHERE `id`=? LIMIT 1", id)
		err = computerRow.Scan(&name, &lastUpdate)
		handleError(err)

		var currentRF int64
		energyDataRow := db.QueryRow("SELECT `RF` FROM `energydata` WHERE `computerID`=? ORDER BY `dateTime` DESC LIMIT 1", id)
		err = energyDataRow.Scan(&currentRF)
		handleError(err)

		computerElement = EnergyComputerResponseElement{
			ID:          id,
			Name:        name,
			MaxRF:       maxRF,
			LastUpdated: lastUpdate,
			CurrentRF:   currentRF,
		}
		computers.ComputerData = append(computers.ComputerData, computerElement)
	}
	computers.Error = false
	return computers, nil
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
