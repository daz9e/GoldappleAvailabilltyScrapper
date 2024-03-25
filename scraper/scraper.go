package scrapper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ga-scraper/db"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func CycleScrapper() {
	for {
		notFoundGoods := CheckLinks()
		fmt.Println(notFoundGoods)
		time.Sleep(time.Minute)
	}
}

func CheckLinks() []string {
	var trueGoodIds []string

	rows, err := db.Database.Query("SELECT good_id FROM links")
	if err != nil {
		fmt.Println("Ошибка запроса к БД:", err)
		return trueGoodIds
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("%s", err)
		}
	}(rows)

	for rows.Next() {
		var goodId string
		if err := rows.Scan(&goodId); err != nil {
			fmt.Println("Ошибка чтения строки:", err)
			continue
		}
		log.Printf("Scrapping item with id: %s", goodId)

		if GetAvailability(goodId) {
			trueGoodIds = append(trueGoodIds, goodId)
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Ошибка при итерации по строкам:", err)
	}

	fmt.Println("ID с доступными ссылками:", trueGoodIds)

	return trueGoodIds
}

func GetAvailability(goodId string) bool {

	url := fmt.Sprintf("https://goldapple.ru/front/api/catalog/product-card?itemId=%s&cityId=2858811e-448a-482e-9863-e03bf06bb5d4&customerGroupId=0", goodId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response ApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding JSON", err)

	}
	fmt.Println(response.Data)

	for _, variant := range response.Data.Variants {
		if variant.InStock {
			return true
		}
	}

	return false
}

type ApiResponse struct {
	Data struct {
		Variants []struct {
			InStock bool `json:"inStock"`
		} `json:"variants"`
	} `json:"data"`
}
