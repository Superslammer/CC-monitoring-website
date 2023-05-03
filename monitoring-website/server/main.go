package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const APIROUTE string = "/api/"
const MYSQL_TIME_FORMAT = "2006-01-02 15:04:05"

// Handle host switching
//type HostSwitch map[string]http.Handler
//
//func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if handler := hs[r.Host]; handler != nil {
//		handler.ServeHTTP(w, r)
//	} else {
//		http.Error(w, "Forbidden", http.StatusForbidden)
//	}
//}

func main() {
	// Setting up httprouter
	website := httprouter.New()
	website.RedirectTrailingSlash = true
	website.HandleOPTIONS = true

	// Make host switch
	//hs := make(HostSwitch)
	//hs["localhost:3000"] = website

	/// Handle api requests
	/// Route: energy-data/
	// Energy data extration
	website.GET(APIROUTE+"energy-data/", extractParameters(withCORS(sendEnergyData)))

	// Energy data insertion
	website.POST(APIROUTE+"energy-data/", extractParameters(withCORS(recieveEnergyData)))

	/// Route: energy-data/:id
	// Get single energy data entry
	website.GET(APIROUTE+"energy-data/:id/", extractParameters(withCORS(sendEnergyData)))

	// Update energy data entry
	website.PATCH(APIROUTE+"energy-data/:id/", extractParameters(withCORS(updateEnergyData)))

	// Remove energy data entry
	website.DELETE(APIROUTE+"energy-data/:id/", extractParameters(withCORS(removeEnergyData)))

	/// Route: energy-computer/
	// Get energy computers
	website.GET(APIROUTE+"energy-computer/", extractParameters(withCORS(sendEnergyComputers)))

	// Assign/Create energy computers
	website.POST(APIROUTE+"energy-computer/", extractParameters(withCORS(createEnergyComputer)))

	/// Route: energy-computer/:id
	// Get a single energy computer
	website.GET(APIROUTE+"energy-computer/:id/", extractParameters(withCORS(sendEnergyComputers)))

	// Update energy computer
	website.PATCH(APIROUTE+"energy-computer/:id/", extractParameters(withCORS(updateEnergyComputer)))

	// Remove/Unassign energy computer
	website.DELETE(APIROUTE+"energy-computer/:id/", extractParameters(withCORS(removeEnergyComputer)))

	// Set CORS
	website.GlobalOPTIONS = http.HandlerFunc(handleCORS)

	// Serve static
	static := httprouter.New()
	static.ServeFiles("/*filepath", http.Dir("./../frontend/dist"))

	website.NotFound = static

	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", website))
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

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func emptyResponse(w http.ResponseWriter) {
	_, err := fmt.Fprint(w, "null")
	handleError(err)
}
