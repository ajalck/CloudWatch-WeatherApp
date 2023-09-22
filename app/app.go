package app

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"reflect"
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

	///
	verifier := struct {
		StatusCod interface{} `json:"cod"`
	}{}
	json.Unmarshal(values, &verifier)
	cod := verifier.StatusCod
	if reflect.ValueOf(cod).Kind() == reflect.Float64 {
		if cod.(float64) != float64(200) {
			return &config.WeatherData{}, errors.New("nod data found")
		}
	} else if reflect.ValueOf(cod).Kind() == reflect.String {
		if cod.(string) != "200" {
			return &config.WeatherData{}, errors.New("nod data found")
		}
	}
	///

	var d config.WeatherData
	err = json.Unmarshal(values, &d)
	if err != nil {
		log.Println("error is: ", err)
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
	sunRise := riseTime.UTC().Add(time.Second * time.Duration(timezoneOffsetSeconds))
	sunset := setTime.UTC().Add(time.Second * time.Duration(timezoneOffsetSeconds))

	// Getting the Country Name

	country, _ := GetCountryName(weatherData.Sys.Country)

	// Setting Background image
	var bgimg string
	if now.After(sunRise) && now.Before(sunset) {
		bgimg = "https://png.pngtree.com/thumb_back/fh260/background/20201012/pngtree-white-cloud-on-blue-sky-weather-background-image_410050.jpg"
	} else if now.After(sunset) || now.Before(sunRise) {
		bgimg = "https://wallpapercave.com/wp/wp9017481.jpg"
	}
	weatherData.Weather[0].Image = weatherImg
	weatherData.Main.Temp = Temp
	weatherData.Main.FeelsLike = TempFeels
	weatherData.Wind.Speed = speed
	weatherData.Now.Time = currentTime
	weatherData.Now.Date = currentDate
	weatherData.Now.Day = currentDay
	weatherData.LocalSunRise = sunRise.Format("03:04 PM")
	weatherData.LocalSunSet = sunset.Format("03:04 PM")
	weatherData.CountryName = country
	weatherData.BGImage = bgimg

	return weatherData
}
