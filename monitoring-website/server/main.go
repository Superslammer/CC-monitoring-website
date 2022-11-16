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
	var numEntries int64 = 20
	var curTime string

	params := make(map[string]string)
	params["dateTime"] = r.URL.Query().Get("dateTime")
	params["numEntries"] = r.URL.Query().Get("numEntries")

	if params["numEntries"] != "" {
		var err error
		numEntries, err = strconv.ParseInt(params["numEntries"], 10, 64)
		checkError(err)
	}

	if params["dateTime"] != "" {
		curTime = params["dateTime"]
	} else {
		curTime = time.Now().Format(MYSQL_TIME_FORMAT)
	}

	data := db.GetEnergyData(numEntries, curTime)
	json.NewEncoder(w).Encode(data)
}

type EnergyData struct {
	Data []Datum `json:"data"`
}

type Datum struct {
	DateTime string `json:"dateTime"`
	RF       int64  `json:"RF"`
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
		db.InsertEnergyData(element.DateTime, element.RF)
	}

	// Return OK or error
	_, err = fmt.Fprintf(w, "Data inserted")
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
