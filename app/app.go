package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"WeatherApp/config"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}

func GetWeatherData(city string) (*config.WeatherData, error) {
	apiConfig, err := config.LoadApiConfig(".apiConfig")
	if err != nil {
		return &config.WeatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return &config.WeatherData{}, err
	}
	defer resp.Body.Close()
	values, _ := io.ReadAll(resp.Body)
	log.Println(string(values))
	var d config.WeatherData
	err = json.Unmarshal(values, &d)
	if err != nil {
		log.Println(err)
		return &config.WeatherData{}, err
	}
	return &d, nil
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
	err = json.Unmarshal(values, &getCountry)
	if err != nil {
		return "", err
	}
	return getCountry[0].Name.Common, nil
}
func SetUpOutput(weatherData *config.WeatherData) *config.WeatherData {
	//Get Weather Image
	getWeatherImage := func(icon string) (image string) {
		image = "https://openweathermap.org/img/wn/" + icon + "@4x.png"
		return
	}
	weatherIcon := weatherData.Weather[0].Icon
	weatherImg := getWeatherImage(weatherIcon)

	//Temperature: Convert Kelvin to degree celcius
	Temp := weatherData.Main.Temp - 273.15
	TempFeels := weatherData.Main.FeelsLike - 273.15

	//Speed: Convert m/s to km/h
	speed := weatherData.Wind.Speed * (18 / 5)

	//Time: Calculate local time of the city

	timezoneOffsetSeconds := weatherData.TimeZoneOffset
	now := time.Now().UTC().Add(time.Second * time.Duration(timezoneOffsetSeconds))

	currentTime := now.Format("03:04 PM")
	currentDate := now.Format("02-01-2006")
	currentDay := now.Weekday()

	// Convert the Unix timestamp to a time.Time value

	riseTime := time.Unix(weatherData.Sys.SunRise, 0)
	setTime := time.Unix(weatherData.Sys.Sunset, 0)
	sunRise := riseTime.UTC().Add(time.Second * time.Duration(timezoneOffsetSeconds)).Format("03:04 PM")
	sunset := setTime.UTC().Add(time.Second * time.Duration(timezoneOffsetSeconds)).Format("03:04 PM")

	// Getting the Country Name

	country, _ := GetCountryName(weatherData.Sys.Country)

	// Setting Background image
	var bgimg string
	if now.Before(setTime) {
		bgimg = "https://png.pngtree.com/thumb_back/fh260/background/20201012/pngtree-white-cloud-on-blue-sky-weather-background-image_410050.jpg"
	} else if now.After(riseTime) {
		bgimg = "https://mdbootstrap.com/img/new/fluid/nature/015.jpg"
	}

	weatherData.Weather[0].Image = weatherImg
	weatherData.Main.Temp = Temp
	weatherData.Main.FeelsLike = TempFeels
	weatherData.Wind.Speed = speed
	weatherData.Now.Time = currentTime
	weatherData.Now.Date = currentDate
	weatherData.Now.Day = currentDay
	weatherData.LocalSunRise = sunRise
	weatherData.LocalSunSet = sunset
	weatherData.CountryName = country
	weatherData.BGImage = bgimg

	return weatherData
}
