package main

import (
	"ga-scraper/db"
	"ga-scraper/tgbot"
	"time"
)

func main() {
	db.InitDB()
	defer db.DatabaseClose()

	tgbot.InitBot()
	go tgbot.ListenBot()
	time.Sleep(10 * time.Minute)
}
