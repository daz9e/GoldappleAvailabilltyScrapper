package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB() *sql.DB {
	fmt.Println("startdb")
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Panic(err)
	}
	if db == nil {
		log.Panic("DB is nil")
	}

	Migrate(db)

	return db
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
