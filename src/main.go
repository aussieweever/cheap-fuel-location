package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"njetch.com/cheap-fuel-location.git/src/models"
)

func main() {
	res, err := http.Get("https://projectzerothree.info/api.php?format=json")
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var queryResult models.QueryResult

	json.Unmarshal(content, &queryResult)
	var priceItems []models.PriceItem

	for _, region := range queryResult.Regions {
		priceItems = append(priceItems, region.Prices...)
	}

	var cheapestPriceItem *models.PriceItem

	for index := range priceItems {
		if priceItems[index].Type == "Diesel" {
			if cheapestPriceItem != nil {
				if cheapestPriceItem.Price > priceItems[index].Price {
					cheapestPriceItem = &priceItems[index]
				}
			} else {
				cheapestPriceItem = &priceItems[index]
			}
		}
	}

	fmt.Printf("The cheapest diesel is in %s %s at %f\n", cheapestPriceItem.Name, cheapestPriceItem.State, cheapestPriceItem.Price)
	fmt.Printf("Lat: %f, Lon: %f\n", cheapestPriceItem.Lat, cheapestPriceItem.Lng)
}
