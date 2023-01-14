package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const APIROUTE string = "/api/"
const MYSQL_TIME_FORMAT = "2006-01-02 15:04:05"

func main() {
	// Setting up httprouter
	r := httprouter.New()
	r.RedirectTrailingSlash = true
	r.HandleOPTIONS = true

	/// Handle api requests
	/// Route: energy-data/
	// Energy data extration
	r.GET(APIROUTE+"energy-data/", extractParameters(sendEnergyData))

	// Energy data insertion
	r.POST(APIROUTE+"energy-data/", extractParameters(recieveEnergyData))

	/// Route: energy-data/:id
	// Get single energy data entry
	r.GET(APIROUTE+"energy-data/:id/", extractParameters(sendEnergyData))

	// Update energy data entry
	r.PATCH(APIROUTE+"energy-data/:id/", extractParameters(updateEnergyData))

	// Remove energy data entry
	r.DELETE(APIROUTE+"energy-data/:id/", extractParameters(removeEnergyData))

	/// Route: energy-computer/
	// Get energy computers
	r.GET(APIROUTE+"energy-computer/", extractParameters(sendEnergyComputers))

	// Assign/Create energy computers
	r.POST(APIROUTE+"energy-computer/", extractParameters(createEnergyComputer))

	/// Route: energy-computer/:id
	// Get a single energy computer
	r.GET(APIROUTE+"energy-computer/:id/", extractParameters(sendEnergyComputers))

	// Update energy computer
	r.PATCH(APIROUTE+"energy-computer/:id/", extractParameters(updateEnergyComputer))

	// Remove/Unassign energy computer
	r.DELETE(APIROUTE+"energy-computer/:id/", extractParameters(removeEnergyComputer))

	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func extractParameters(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		params := r.URL.Query()
		for key, val := range params {
			par := httprouter.Param{Key: key, Value: val[0]}
			ps = append(ps, par)
		}
		fn(w, r, ps)
	}
}

/*func withAPIKey(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, http.StatusUnauthorized, "Invalid API key")
			return
		}
		ps = append(ps, httprouter.Param{Key: "key", Value: key})
		fn(w, r, ps)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}*/

type ErrorResponse struct {
	Error string `json:"error"`
}

/*func respondErr(w http.ResponseWriter, status int, message string) {
	var response ErrorResponse
	response.Error = message
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}*/

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func emptyResponse(w http.ResponseWriter) {
	_, err := fmt.Fprint(w, "{}")
	handleError(err)
}
