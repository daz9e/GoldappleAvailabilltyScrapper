package tgbot

import (
	"ga-scraper/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
)

var Bot *tgbotapi.BotAPI

func InitBot() {
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

	log.Printf("Authorized on account %s", bot.Self.UserName)

	Bot = bot
}

func ListenBot() {
	InitBot()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		db.CreateUser(update.Message.From)

		log.Printf("User id %s:", update.Message.From.ID)

		re := regexp.MustCompile(`https://goldapple.ru/(\d+)-`)
		matches := re.FindStringSubmatch(update.Message.Text)
		if len(matches) > 1 {
			productID := matches[1]
			db.CreateLink(productID, update.Message.Text)
			db.CreateManyToManyForLinks(update.Message.From, productID)
			//url := fmt.Sprintf("https://goldapple.ru/front/api/catalog/product-card?itemId=%s&cityId=2858811e-448a-482e-9863-e03bf06bb5d4&customerGroupId=0", productID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Товар добавлен")
			_, err := Bot.Send(msg)
			if err != nil {
				return
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сылка не подходит")
			_, err := Bot.Send(msg)
			if err != nil {
				return
			}

		}
	}
}
