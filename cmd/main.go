package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ajalck/WeatherApp/app"
)

func main() {
	http.HandleFunc("/hello",app.Hello)

	http.HandleFunc("/weather/",func(w http.ResponseWriter, r *http.Request) {
		city :=strings.SplitN(r.URL.Path,"/",3)[2]
		data,err:=app.Query(city)
		if err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type","application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})


	http.ListenAndServe(":8080",nil)
}