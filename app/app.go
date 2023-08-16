package app

import (
	"encoding/json"
	"net/http"

	"WeatherApp/config"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}

func Query(city string) (config.WeatherData, error) {
	apiConfig, err := config.LoadApiConfig(".apiConfig")
	if err != nil {
		return config.WeatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return config.WeatherData{}, err
	}
	defer resp.Body.Close()

	var d config.WeatherData
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return config.WeatherData{}, err
	}
	return d, err
}
