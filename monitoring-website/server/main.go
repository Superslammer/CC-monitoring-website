package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const APIROUTE string = "/"
const MYSQL_TIME_FORMAT = "2006-01-02 15:04:05"

// Handle host switching
type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func main() {
	// Setting up httprouter
	api := httprouter.New()
	api.RedirectTrailingSlash = true
	api.HandleOPTIONS = true

	// Make host switch
	hs := make(HostSwitch)
	hs["api.localhost:3000"] = api

	/// Handle api requests
	/// Route: energy-data/
	// Energy data extration
	api.GET(APIROUTE+"energy-data/", extractParameters(withCORS(sendEnergyData)))

	// Energy data insertion
	api.POST(APIROUTE+"energy-data/", extractParameters(withCORS(recieveEnergyData)))

	/// Route: energy-data/:id
	// Get single energy data entry
	api.GET(APIROUTE+"energy-data/:id/", extractParameters(withCORS(sendEnergyData)))

	// Update energy data entry
	api.PATCH(APIROUTE+"energy-data/:id/", extractParameters(withCORS(updateEnergyData)))

	// Remove energy data entry
	api.DELETE(APIROUTE+"energy-data/:id/", extractParameters(withCORS(removeEnergyData)))

	/// Route: energy-computer/
	// Get energy computers
	api.GET(APIROUTE+"energy-computer/", extractParameters(withCORS(sendEnergyComputers)))

	// Assign/Create energy computers
	api.POST(APIROUTE+"energy-computer/", extractParameters(withCORS(createEnergyComputer)))

	/// Route: energy-computer/:id
	// Get a single energy computer
	api.GET(APIROUTE+"energy-computer/:id/", extractParameters(withCORS(sendEnergyComputers)))

	// Update energy computer
	api.PATCH(APIROUTE+"energy-computer/:id/", extractParameters(withCORS(updateEnergyComputer)))

	// Remove/Unassign energy computer
	api.DELETE(APIROUTE+"energy-computer/:id/", extractParameters(withCORS(removeEnergyComputer)))

	// Set CORS
	api.GlobalOPTIONS = http.HandlerFunc(handleCORS)

	// Serve website
	website := httprouter.New()
	website.ServeFiles("/*filepath", http.Dir("./../frontend/dist"))

	hs["localhost:3000"] = website

	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", hs))
}

func handleCORS(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Access-Control-Allow-Origin", "*")
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

func withCORS(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
