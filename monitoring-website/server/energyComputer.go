package main

import (
	"api-design/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func sendEnergyComputers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Database connection
	db := database.ConnectDB()

	// Get numComputers as number
	var numComputers int
	if numEntriesPar := ps.ByName("numComputers"); numEntriesPar != "" {
		var err error
		numComputers, err = strconv.Atoi(numEntriesPar)
		handleError(err)
	} else {
		numComputers = 20
	}

	var data []database.EnergyComputer
	var err error
	parID := ps.ByName("id")
	if parID == "" {
		data, err = db.GetEnergyComputers(-1, numComputers)
	} else if parID != "" {
		var id int
		id, err = strconv.Atoi(parID)
		handleError(err)
		data, err = db.GetEnergyComputers(id, 1)
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

	json.NewEncoder(w).Encode(&data)
}

type EnergyComputers []EnergyComputer

type EnergyComputer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	MaxRF     int64  `json:"maxRF"`
	CurrentRF int64  `json:"currentRF"`
}

func createEnergyComputer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Datbase connection
	db := database.ConnectDB()

	// Decoding recived data
	var decoded EnergyComputers
	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		emptyResponse(w)
		fmt.Println(err)
		return
	}

	// Create/Assign energy computers
	for _, computer := range decoded {
		err := db.CreateOrAssignEnergyComputer(computer.ID, computer.Name, computer.MaxRF, computer.CurrentRF)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			emptyResponse(w)
			fmt.Println(err)
			return
		}
	}
}

type EnergyComputerUpdate struct {
	Name  string `json:"name"`
	MaxRF int64  `json:"maxRF"`
}

func updateEnergyComputer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Datbase connection
	db := database.ConnectDB()

	// Decode the recieved data
	var decoded EnergyComputerUpdate
	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		emptyResponse(w)
		fmt.Println(err)
		return
	}

	parID := ps.ByName("id")
	var id int
	id, err = strconv.Atoi(parID)
	handleError(err)

	err = db.UpdateEnergyComputer(id, decoded.Name, decoded.MaxRF)
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

func removeEnergyComputer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Database connection
	db := database.ConnectDB()

	parID := ps.ByName("id")
	id, err := strconv.Atoi(parID)
	handleError(err)

	err = db.RemoveEnergyComputer(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emptyResponse(w)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	emptyResponse(w)
}
