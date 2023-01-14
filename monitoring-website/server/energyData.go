package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"api-design/database"

	"github.com/julienschmidt/httprouter"
)

func sendEnergyData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.ConnectDB()
	defer db.Close()

	// Getting the datetime parameter
	var curTime string
	if timePar := ps.ByName("dateTime"); timePar != "" {
		curTime = timePar
	} else {
		curTime = time.Now().Format(MYSQL_TIME_FORMAT)
	}

	// Getting the number of entries to return
	var numEntries int
	if numEntriesPar := ps.ByName("numEntries"); numEntriesPar != "" {
		var err error
		numEntries, err = strconv.Atoi(numEntriesPar)
		handleError(err)
	} else {
		numEntries = 20
	}

	// Getting wether or not the computerID has been set
	var computerID int
	if computerIDPar := ps.ByName("computerID"); computerIDPar != "" {
		var err error
		computerID, err = strconv.Atoi(computerIDPar)
		handleError(err)
	} else {
		computerID = -1
	}

	var err error
	var data []database.DataPoint
	if strID := ps.ByName("id"); strID != "" {
		var id int
		id, err = strconv.Atoi(strID)
		handleError(err)
		data, err = db.GetEnergyData(id, computerID, numEntries, curTime)
	} else {
		data, err = db.GetEnergyData(-1, computerID, numEntries, curTime)
	}

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		emptyResponse(w)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emptyResponse(w)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

type EnergyData []EnergyDatum

type EnergyDatum struct {
	DateTime   string `json:"dateTime"`
	RF         int64  `json:"RF"`
	ComputerID int    `json:"computerID"`
}

func recieveEnergyData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.ConnectDB()
	defer db.Close()

	var decoded EnergyData
	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		emptyResponse(w)
		fmt.Println(err)
		return
	}

	// Insert data into database
	for _, element := range decoded {
		err = db.InsertEnergyData(element.DateTime, element.RF, element.ComputerID)
		handleError(err)
	}

	// Return OK or error
	w.WriteHeader(http.StatusOK)
	emptyResponse(w)
}

func updateEnergyData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.ConnectDB()
	defer db.Close()

	var decoded EnergyDatum
	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		emptyResponse(w)
		return
	}

	parID := ps.ByName("id")
	id, err := strconv.Atoi(parID)
	handleError(err)

	err = db.UpdateEnergyData(id, decoded.DateTime, decoded.RF, decoded.ComputerID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		emptyResponse(w)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emptyResponse(w)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	emptyResponse(w)
}

func removeEnergyData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.ConnectDB()
	defer db.Close()

	parID := ps.ByName("id")
	id, err := strconv.Atoi(parID)
	handleError(err)

	err = db.RemoveEnergyData(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		emptyResponse(w)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emptyResponse(w)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	emptyResponse(w)
}
