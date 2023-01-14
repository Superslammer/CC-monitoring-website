package database

import (
	"database/sql"
	"errors"
	"time"
)

// Custom errors
var ErrAlreadyAssignedEnergy = errors.New("computer already assigned as energy computer")

type EnergyComputer struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MaxRF       int64  `json:"maxRF"`
	LastUpdated string `json:"lastUpdated"`
	CurrentRF   int64  `json:"currentRF"`
}

func (db *Database) GetEnergyComputers(id int, numComputers int) ([]EnergyComputer, error) {
	if id == -1 {
		// Getting the requested energy computers
		energyComputersDBRes, err := db.db.Query("SELECT `id`, `maxRF` FROM `energycomputers` LIMIT ?", numComputers)
		if err != nil {
			nilResponse := []EnergyComputer{}
			return nilResponse, err
		}

		var energyComputers []EnergyComputer
		for energyComputersDBRes.Next() {
			var energyComputer EnergyComputer
			// Getting id and max RF
			err := energyComputersDBRes.Scan(&energyComputer.ID, &energyComputer.MaxRF)
			if err != nil {
				nilResponse := []EnergyComputer{}
				return nilResponse, err
			}

			// Getting name and last update
			computerDBRes := db.db.QueryRow("SELECT `name`, `lastUpdate` FROM `computers` WHERE `id`=?", &energyComputer.ID)
			err = computerDBRes.Scan(&energyComputer.Name, &energyComputer.LastUpdated)
			if err != nil {
				nilResponse := []EnergyComputer{}
				return nilResponse, err
			}

			// Getting current RF
			curRFRes := db.db.QueryRow("SELECT `RF` FROM `energydata` WHERE `computerID`=? ORDER BY `dateTime` DESC LIMIT 1", energyComputer.ID)
			err = curRFRes.Scan(&energyComputer.CurrentRF)
			if err != nil {
				nilResponse := []EnergyComputer{}
				return nilResponse, err
			}

			energyComputers = append(energyComputers, energyComputer)
		}

		return energyComputers, nil
	} else {
		var energyComputer EnergyComputer
		// Getting id and max RF
		energyComputerRes := db.db.QueryRow("SELECT `id`, `maxRF` FROM `energycomputers` WHERE `id`=?", id)
		err := energyComputerRes.Scan(&energyComputer.ID, &energyComputer.MaxRF)
		if err != nil {
			nilResponse := []EnergyComputer{}
			return nilResponse, err
		}

		// Getting name and last update
		computerRes := db.db.QueryRow("SELECT `name`, `lastUpdate` FROM `computers` WHERE `id`=?", energyComputer.ID)
		err = computerRes.Scan(&energyComputer.Name, &energyComputer.LastUpdated)
		if err != nil {
			nilResponse := []EnergyComputer{}
			return nilResponse, err
		}

		// Getting current RF
		curRFRes := db.db.QueryRow("SELECT `RF` FROM `energydata` WHERE `computerID`=? ORDER BY `dateTime` DESC LIMIT 1", energyComputer.ID)
		err = curRFRes.Scan(&energyComputer.CurrentRF)
		if err != nil {
			nilResponse := []EnergyComputer{}
			return nilResponse, err
		}

		var energyComputers []EnergyComputer
		energyComputers = append(energyComputers, energyComputer)
		return energyComputers, nil
	}
}

func (db *Database) CreateOrAssignEnergyComputer(computerID int, name string, maxRF int64, currentRF int64) error {
	if computerID == -1 {
		// Computer not regisered, register computer
		var err error
		computerID, err = db.RegisterComputer(ComputerTypes.EnergyComputer, name)
		if err != nil {
			return err
		}
	} else {
		// Check if computer is registered
		compTest := db.db.QueryRow("SELECT `id` FROM `computers` WHERE `id`=?", computerID)
		var tempID int
		err := compTest.Scan(&tempID)
		if err == sql.ErrNoRows {
			return ErrComputerNotRegistered
		} else if err != nil {
			return err
		}
	}

	// Check if computer is already assigned as energy computer
	eCompTest := db.db.QueryRow("SELECT `id` FROM `energycomputers` WHERE `id`=?", computerID)
	var tempID int
	err := eCompTest.Scan(&tempID)
	if err == sql.ErrNoRows {
		// Assign computer
		_, err := db.db.Exec("INSERT INTO `energycomputers` (`id`, `maxRF`) VALUES(?, ?)", computerID, maxRF)
		handleError(err)
	} else if err != nil {
		return err
	} else {
		return ErrAlreadyAssignedEnergy
	}

	// Insert current RF
	curTime := time.Now().Format(MYSQL_TIME_FORMAT)
	err = db.InsertEnergyData(curTime, currentRF, computerID)
	return err
}

func (db *Database) UpdateEnergyComputer(id int, name string, maxRF int64) error {
	// Check if energy computer exsist
	test := db.db.QueryRow("SELECT `id` FROM `energycomputers` WHERE `id`=?", id)
	var tempID int
	err := test.Scan(&tempID)
	if err != nil {
		return err
	}

	// Update info
	_, err = db.db.Exec("UPDATE `energycomputers` SET `maxRF`=? WHERE `id`=?", maxRF, id)
	if err != nil {
		return err
	}
	_, err = db.db.Exec("UPDATE `computers` SET `name`=? WHERE `id`=?", name, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveEnergyComputer(id int) error {
	// Check if energy computer exsist
	test := db.db.QueryRow("SELECT `id` FROM `energycomputers` WHERE `id`=?", id)
	var tempID int
	err := test.Scan(&tempID)
	if err != nil {
		return err
	}

	// Remove computer from "energycomputers" table
	_, err = db.db.Exec("DELETE FROM `energycomputers` WHERE `id`=?", id)
	if err != nil {
		return err
	}

	// Assign computer the "Not Assigned" computer type
	_, err = db.db.Exec("UPDATE `computers` SET `typeID`=? WHERE `id`=?", ComputerTypes.NotAssigned, id)
	return err
}
