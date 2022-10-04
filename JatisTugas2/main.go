package main

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
)

type apiConfigData struct {
	OpenWeatherApiKey string `json: "OpenWeatherApiKey"`
}

type weatherData struct {
	Name string `json: "nama"`
	Main struct {
		Kelvin float64 `json: "temp"`
	} `json: "main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func query(lat string, lon string)(weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + lon + "&appid=" + apiConfig.OpenWeatherApiKey)

	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData
	if err != json.NewDecoder(resp.Body).Decode(&d) {
		return weatherData{}, err
	}
	return d, nil
}

func main() {
	http.HandleFunc("/weather", 
	func(w http.ResponseWriter, r *http.Request) {
		lat := strings.SplitN(r.URL.Path, "/", 3)[2]
		lon := strings.SplitN(r.URL.Path, "/", 4)[3]
		data, err := query(lat, lon)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset = utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}