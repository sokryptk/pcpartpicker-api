package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"pcpartpicker-api/api/endpoints"
	"pcpartpicker-api/scraper"
)




func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":6920"
}

func main() {
	defer scraper.Instance.Quit()
	defer scraper.Service.Stop()

	router := mux.NewRouter()

	router.HandleFunc("/guides", endpoints.GetBuildGuides).Methods("GET")
	router.HandleFunc("/gdetails", endpoints.GetGuideDetails)
	router.HandleFunc("/parts", endpoints.GetPartsList)
	router.HandleFunc("/completedBuilds", endpoints.GetCompletedBuilds).Methods("POST")


	_ = http.ListenAndServe(getPort(), router)
}
