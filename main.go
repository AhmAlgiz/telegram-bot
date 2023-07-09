package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var WeatherApi string

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

func GetWeather(loc string) string {
	req, err := http.Get(fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", WeatherApi, url.QueryEscape(loc)))
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

func main() {
	//telegram bot token is the 1st parameter, wheaterapi.com API key is 2nd
	tgToken := os.Args[1]
	WeatherApi = os.Args[2]
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, GetWeather(update.Message.Text))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
