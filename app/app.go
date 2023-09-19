package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"WeatherApp/config"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}

func GetWeatherData(city string) (config.WeatherData, error) {
	apiConfig, err := config.LoadApiConfig(".apiConfig")
	if err != nil {
		return config.WeatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return config.WeatherData{}, err
	}
	defer resp.Body.Close()
	values, _ := io.ReadAll(resp.Body)
	log.Println(string(values))
	var d config.WeatherData
	err = json.Unmarshal(values, &d)
	if err != nil {
		log.Println(err)
		return config.WeatherData{}, err
	}
	return d, nil
}
func GetCountryName(code string) (string, error) {
	resp, err := http.Get("https://restcountries.com/v3.1/alpha/" + code)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	values, _ := io.ReadAll(resp.Body)
	getCountry := []struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
	}{}
	json.Unmarshal(values,&getCountry)
	fmt.Println(getCountry)
	return getCountry[0].Name.Common, nil
}
//func SetUpOutput(weatherData config.WeatherData)