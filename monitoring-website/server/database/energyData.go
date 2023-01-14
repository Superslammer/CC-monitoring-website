package database

import (
	"database/sql"
	"time"
)

type DataPoint struct {
	ID         int    `json:"ID"`
	DateTime   string `json:"dateTime"`
	RF         int64  `json:"RF"`
	ComputerID int    `json:"computerID"`
}

func (db *Database) GetEnergyData(id int, computerID int, numEntries int, startDateTime string) ([]DataPoint, error) {
	var err error

	dateTime := db.db.QueryRow("SELECT `dateTime` FROM `energydata` "+
		"ORDER BY ABS(TIMESTAMPDIFF(SECOND, `datetime`, ?)) "+
		"LIMIT 1", startDateTime)

	var closestDateTime string
	err = dateTime.Scan(&closestDateTime)
	if err != nil {
		nilData := []DataPoint{}
		return nilData, err
	}

	var dataRows *sql.Rows
	if id == -1 {
		if computerID == -1 {
			dataRows, err = db.db.Query("SELECT `dateTime`, `RF`, `computerID`, `id` "+
				"FROM `energydata` "+
				"WHERE `dateTime` <= ? "+
				"ORDER BY `dateTime` DESC "+
				"LIMIT ?",
				closestDateTime,
				numEntries)
		} else {
			dataRows, err = db.db.Query("SELECT `dateTime`, `RF`, `computerID`, `id` "+
				"FROM `energydata` "+
				"WHERE `computerID` = ? AND `dateTime` <= ? "+
				"ORDER BY `dateTime` DESC "+
				"LIMIT ?",
				computerID, closestDateTime, numEntries)
		}
		defer dataRows.Close()
	} else {
		dataRows, err = db.db.Query("SELECT `dateTime`, `RF`, `computerID`, `id` "+
			"FROM `energydata` "+
			"WHERE `id` = ? ",
			id)
	}

	if err != nil {
		nilData := []DataPoint{}
		return nilData, err
	}

	var dataSet []DataPoint
	for dataRows.Next() {
		var dataPoint DataPoint
		err = dataRows.Scan(&dataPoint.DateTime, &dataPoint.RF, &dataPoint.ComputerID, &dataPoint.ID)
		handleError(err)
		dataSet = append(dataSet, dataPoint)
	}
	return dataSet, nil
}

func (db *Database) InsertEnergyData(dateTime string, rf int64, computerID int) error {
	_, err := db.db.Exec("INSERT INTO `energydata` (`dateTime`, `RF`, `computerID`) VALUES(?, ?, ?)", dateTime, rf, computerID)
	if err != nil {
		return err
	}

	// Update `lastUpdated` in energyComputers
	curTime := time.Now().Format(MYSQL_TIME_FORMAT)
	_, err = db.db.Exec("UPDATE `computers` SET `lastUpdate`=? WHERE `id`=?", curTime, computerID)
	return err
}

func (db *Database) UpdateEnergyData(id int, datetime string, rf int64, computerId int) error {
	// Check if entry exsist
	test := db.db.QueryRow("SELECT `id` FROM `energydata` WHERE `id`=?", id)
	var tempID int
	err := test.Scan(&tempID)
	if err != nil {
		return err
	}

	_, err = db.db.Exec("UPDATE `energydata` "+
		"SET `dateTime`=?, `RF`=?, `computerID`=? "+
		"WHERE `id`=?",
		datetime, rf, computerId, id)
	return err
}

func (db *Database) RemoveEnergyData(id int) error {
	// Check if entry exsist
	test := db.db.QueryRow("SELECT `id` FROM `energydata` WHERE `id`=?", id)
	var tempID int
	err := test.Scan(&tempID)
	if err != nil {
		return err
	}

	_, err = db.db.Exec("DELETE FROM `energydata` WHERE `id`=?", id)
	return err
}
