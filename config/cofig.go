package config

import (
	"encoding/json"
	"time"

	"os"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type WeatherData struct {
	Name    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Image       string `json:"image"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Sys struct {
		Country string `json:"country"`
		SunRise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	TimeZoneOffset int64 `json:"timezone"`

	//
	Now struct {
		Time string       `json:"time"`
		Date string       `json:"date"`
		Day  time.Weekday `json:"day"`
	} `json:"now"`
	LocalSunRise string `json:"localsunrise"`
	LocalSunSet  string `json:"localsunset"`
	CountryName  string `json:"countryname"`
	BGImage      string `json:"bgimage"`
}

func LoadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	if err = json.Unmarshal(bytes, &c); err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}
