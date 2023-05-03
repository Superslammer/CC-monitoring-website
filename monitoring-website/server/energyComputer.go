package main

import (
	"api-design/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func sendEnergyComputers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Database connection
	db := database.ConnectDB()
	defer db.Close()

	// Get numComputers as number
	var numComputers int
	if numEntriesPar := ps.ByName("numComputers"); numEntriesPar != "" {
		var err error
		numComputers, err = strconv.Atoi(numEntriesPar)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			emptyResponse(w)
		}
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
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			emptyResponse(w)
		}
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
		log.Panic(err)
		return
	}

	fmt.Println(data)

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
	defer db.Close()

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
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w)
}

type EnergyComputerUpdate struct {
	Name  string `json:"name"`
	MaxRF int64  `json:"maxRF"`
}

func updateEnergyComputer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Datbase connection
	db := database.ConnectDB()
	defer db.Close()

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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		emptyResponse(w)
	}

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
	defer db.Close()

	parID := ps.ByName("id")
	id, err := strconv.Atoi(parID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		emptyResponse(w)
	}

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
