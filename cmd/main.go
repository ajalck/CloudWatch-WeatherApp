package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"WeatherApp/app"
)

func main() {
	http.HandleFunc("/hello", app.Hello)
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		htmlPage := "./app/index.html"
		http.ServeFile(w, r, htmlPage)
	})
	http.HandleFunc("/weather/:city", func(w http.ResponseWriter, r *http.Request) {
		City := r.URL.Query().Get("city")
		data, err := app.GetWeatherData(City)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		weatherData := app.SetUpOutput(data)

		jsonData, err := json.Marshal(weatherData)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprint(w, string(jsonData))
	})
	log.Println("Server Sarting at: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Failed to start server")
	}
}
