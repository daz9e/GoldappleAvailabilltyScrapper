package scrapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetAvailability(url string) string {

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return "Error in request"
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "Error in client"
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "Error in ioutil"
	}

	var response ApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding JSON", err)
		return "Error in JSON decoding"
	}
	fmt.Println(response.Data)

	for _, variant := range response.Data.Variants {
		if variant.InStock {
			return "Product is available"
		}
	}

	return "Product is not available"
}

type ApiResponse struct {
	Data struct {
		Variants []struct {
			InStock bool `json:"inStock"`
		} `json:"variants"`
	} `json:"data"`
}
