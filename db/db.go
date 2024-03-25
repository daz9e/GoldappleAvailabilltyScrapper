package db

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	Database *sql.DB
)

func InitDB() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Panic(err)
	}
	if db == nil {
		log.Panic("DB is nil")
	}

	Migrate(db)

	Database = db
}

func DatabaseClose() {
	err := Database.Close()
	if err != nil {
		log.Printf("Error with closing database: %s", err)
	}
}

func Migrate(db *sql.DB) {
	sqlCreateQuery := `
    CREATE TABLE IF NOT EXISTS users(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        telegram_id INTEGER NOT NULL UNIQUE,
        username TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS links(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        good_id INTEGER NOT NULL UNIQUE,
        url TEXT NOT NULL UNIQUE
    );
    CREATE TABLE IF NOT EXISTS user_links(
        user_id INTEGER NOT NULL,
        link_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (link_id) REFERENCES links(id),
        PRIMARY KEY (user_id, link_id)
    );`

	_, err := db.Exec(sqlCreateQuery)
	if err != nil {
		log.Panic(err)
	}
}

func FindUsersByGoods(userIds []string) {
}

func CreateUser(user *tgbotapi.User) {
	_, err := Database.Exec("INSERT INTO users (telegram_id, username) VALUES (?, ?)", user.ID, user.UserName)
	if err != nil {
		log.Printf("Error with creating new user: %s", err)
	}

}

func CreateLink(goodId string, url string) {
	_, err := Database.Exec("INSERT INTO links (good_id, url) VALUES (?, ?)", goodId, url)
	if err != nil {
		log.Printf("Error with creating new link: %s", err)
	}
}

func CreateManyToManyForLinks(user *tgbotapi.User, goodId string) {
	_, err := Database.Exec("INSERT INTO user_links (user_id, link_id) VALUES (?, ?)", user.ID, goodId)
	if err != nil {
		log.Printf("Error with creating new user_link: %s", err)
	}

}
func AddGoodForUser(userId string) {

}
