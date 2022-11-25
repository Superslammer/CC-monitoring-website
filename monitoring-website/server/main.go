package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	db "github.com/Superslammer/CC-monitoring-website/database"
	"github.com/gorilla/mux"
)

const MYSQL_TIME_FORMAT string = "2006-01-02 15:04:05"

func main() {
	// Connect to database
	db.ConnectDB()

	// Setting up gorilla mux
	r := mux.NewRouter().StrictSlash(true)

	// Handle api requests
	api := r.PathPrefix("/api/").Subrouter()

	// Respond that the api is online
	api.HandleFunc("/", isOnline).Methods("GET")

	// Register computer as energy computer
	api.HandleFunc("/energyComputer", registerEnergyComputer).Methods("POST")

	// Get computers registered as energy
	api.HandleFunc("/energyComputer", sendEnergyComputers).Methods("GET")

	// Update energy computer info

	// Energy data extration
	api.HandleFunc("/energyData", sendEnergyData).Methods("GET")

	// Handle energy data insertion
	api.HandleFunc("/energyData", recieveEnergyData).Methods("POST")

	// Serve website
	fs := http.FileServer(http.Dir("./../frontend/dist"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	fmt.Println("Server listening on port 3000")
	log.Panic(
		http.ListenAndServe(":3000", r),
	)
}

// Handling sending energy data to computercraft
func sendEnergyData(w http.ResponseWriter, r *http.Request) {
	var numEntries int = 20
	var curTime string
	var computerID int

	params := make(map[string]string)
	params["dateTime"] = r.URL.Query().Get("dateTime")
	params["numEntries"] = r.URL.Query().Get("numEntries")
	params["computerID"] = r.URL.Query().Get("computerID")

	if params["numEntries"] != "" {
		var err error
		numEntries, err = strconv.Atoi(params["numEntries"])
		handleError(err)
	}

	if params["dateTime"] != "" {
		curTime = params["dateTime"]
	} else {
		curTime = time.Now().Format(MYSQL_TIME_FORMAT)
	}

	if params["computerID"] != "" {
		var err error
		computerID, err = strconv.Atoi(params["computerID"])
		handleError(err)
	} else {
		computerID = -1
	}

	data, err := db.GetEnergyData(computerID, numEntries, curTime)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(&data)
}

type EnergyData struct {
	Data []Datum `json:"data"`
}

type Datum struct {
	DateTime   string `json:"dateTime"`
	RF         int64  `json:"RF"`
	ComputerID int    `json:"computerID"`
}

// Handling recieving energy data from computercraft
func recieveEnergyData(w http.ResponseWriter, r *http.Request) {
	// Decode data
	var decoded EnergyData

	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert data into database
	for _, element := range decoded.Data {
		db.InsertEnergyData(element.DateTime, element.RF, element.ComputerID)
	}

	// Return OK or error
	_, err = fmt.Fprintf(w, "Data inserted")
	handleError(err)
}

type EnergyComputerRegistrationRequest struct {
	ComputerID int    `json:"computerID"`
	MaxEnergy  int64  `json:"maxEnergy"`
	Name       string `json:"name"`
}

type ComputerRegistrationResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
}

func registerEnergyComputer(w http.ResponseWriter, r *http.Request) {
	// Decode data
	var decoded EnergyComputerRegistrationRequest

	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var response ComputerRegistrationResponse
	err = db.CreateEnergyComputerEntry(decoded.ComputerID, decoded.MaxEnergy, decoded.Name)
	if err != nil {
		response = ComputerRegistrationResponse{
			Error: true,
			Msg:   err.Error(),
		}
	} else {
		response = ComputerRegistrationResponse{
			Error: false,
			Msg:   "Computer registered as energy computer",
		}
	}
	json.NewEncoder(w).Encode(&response)
}

func sendEnergyComputers(w http.ResponseWriter, r *http.Request) {
	var numComputers int = 10

	params := make(map[string]string)
	params["numComputers"] = r.URL.Query().Get("numComputers")
	params["computerID"] = r.URL.Query().Get("computerID")

	compId, isNumCompId := strconv.Atoi(params["computerID"])
	newNumComputers, isNumNumComputers := strconv.Atoi(params["numComputers"])
	if params["computerID"] != "" && isNumCompId == nil {
		// Get specific computer data
		computers, err := db.GetEnergyComputers(compId, 1)
		handleError(err)
		json.NewEncoder(w).Encode(computers)
		return
	} else if params["computerID"] != "" && isNumCompId != nil {
		parameterError("computerID", "int/number", w)
		return
	} else if params["numComputers"] != "" && isNumNumComputers == nil {
		computers, err := db.GetEnergyComputers(-1, newNumComputers)
		handleError(err)
		json.NewEncoder(w).Encode(computers)
		return
	} else if params["numComputers"] != "" && isNumNumComputers != nil {
		parameterError("numComputers", "int/number", w)
		return
	}

	computers, err := db.GetEnergyComputers(-1, numComputers)
	handleError(err)
	json.NewEncoder(w).Encode(computers)
}

func parameterError(parameter string, paraType string, w http.ResponseWriter) {
	fmt.Fprintf(w, `{"error": true, "msg": "Invalid value in parameter: %s. Must be of type: %s"}`, parameter, paraType)
}

func isOnline(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Online")
}

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
