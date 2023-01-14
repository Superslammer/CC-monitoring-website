package database

import (
	"errors"
	"time"
)

// Custom errors
var ErrComputerNotRegistered = errors.New("desired computer not registered")

type computerTypes struct {
	NotAssigned    int
	EnergyComputer int
}

var ComputerTypes computerTypes = computerTypes{
	NotAssigned:    0,
	EnergyComputer: 1,
}

func (db *Database) RegisterComputer(typeID int, name string) (int, error) {
	curTime := time.Now().Format(MYSQL_TIME_FORMAT)
	result, err := db.db.Exec("INSERT INTO `computers` (`typeID`, `name`, `lastUpdate`) VALUES(?, ?, ?)", typeID, name, curTime)
	if err != nil {
		return -1, err
	}
	var id int64
	id, err = result.LastInsertId()
	return int(id), err
}
