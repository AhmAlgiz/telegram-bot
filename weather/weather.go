package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Location struct {
	Name    string
	Country string
}
type Current struct {
	Temp_c    float32
	Condition Condition
	Wind_kph  float32
	Wind_dir  string
}

type Condition struct {
	Text string
}

type Weather struct {
	Location Location
	Current  Current
}

func Get(weatherAPI, loc string) string {
	req, err := http.Get(fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherAPI, url.QueryEscape(loc)))
	if err != nil {
		return "Weather service error. Please, try later."
	}
	defer req.Body.Close()

	var w Weather
	err = json.NewDecoder(req.Body).Decode(&w)
	if err != nil {
		return "Request parsing error"
	}

	resp := fmt.Sprintf("Today in %s, %s\nTemperature: %.1f°C\nWind Speed: %.1fmps (%.1fkph), %s\nCondition: %s\n",
		w.Location.Name, w.Location.Country, w.Current.Temp_c, w.Current.Wind_kph*1000/3600, w.Current.Wind_kph, w.Current.Wind_dir, w.Current.Condition.Text)

	return resp
}
