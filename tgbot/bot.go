package tgbot

import (
	"fmt"
	scraper "ga-scraper/scraper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
)

func InitBot() *tgbotapi.BotAPI {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file: ", err)
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Panic("TELEGRAM_BOT_TOKEN is required")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

func ListenBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//TODO: Change logic
		re := regexp.MustCompile(`https://goldapple.ru/(\d+)-`)
		matches := re.FindStringSubmatch(update.Message.Text)
		if len(matches) > 1 {
			productID := matches[1]

			url := fmt.Sprintf("https://goldapple.ru/front/api/catalog/product-card?itemId=%s&cityId=2858811e-448a-482e-9863-e03bf06bb5d4&customerGroupId=0", productID)

			responseBody := scraper.GetAvailability(url)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseBody)
			bot.Send(msg)
		}
	}
}
