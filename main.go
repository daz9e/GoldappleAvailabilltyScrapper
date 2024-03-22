package main

import (
	db "ga-scraper/db"
	"ga-scraper/tgbot"
)

func main() {
	database := db.InitDB()
	err := database.Close()
	if err != nil {
		return
	}
	botInstance := tgbot.InitBot()
	go tgbot.ListenBot(botInstance)
}
