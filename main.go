package main

import (
	"ga-scraper/db"
	scrapper "ga-scraper/scraper"
	"ga-scraper/tgbot"
)

func main() {
	db.InitDB()
	defer db.DatabaseClose()

	go tgbot.ListenBot()
	scrapper.CycleScrapper()
}
