package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"njetch.com/cheap-fuel-location.git/src/models"
	xmlmodel "njetch.com/cheap-fuel-location.git/src/models/xml"
)

var xmlHeader = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	fuelType := os.Getenv("FUEL_TYPE")
	fmt.Printf("Fuel type: %s\n", fuelType)

	configurationFile := os.Getenv("FILE_PATH")
	fmt.Println("Configuration file: " + configurationFile)

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
		if priceItems[index].Type == fuelType {
			if cheapestPriceItem != nil {
				if cheapestPriceItem.Price > priceItems[index].Price {
					cheapestPriceItem = &priceItems[index]
				}
			} else {
				cheapestPriceItem = &priceItems[index]
			}
		}
	}

	fmt.Println("")

	fmt.Printf("The cheapest diesel is in %s %s at %f\n", cheapestPriceItem.Name, cheapestPriceItem.State, cheapestPriceItem.Price)
	fmt.Printf("Lat: %f, Lon: %f\n", cheapestPriceItem.Lat, cheapestPriceItem.Lng)

	xmlFile, err := os.Open(configurationFile)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	xmlContent, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Println(err)
	}
	var gpxTag xmlmodel.GpxTag

	xml.Unmarshal(xmlContent, &gpxTag)

	gpxTag.Position.Lat = fmt.Sprintf("%f", cheapestPriceItem.Lat)
	gpxTag.Position.Lon = fmt.Sprintf("%f", cheapestPriceItem.Lng)

	updatedContent, err := xml.MarshalIndent(gpxTag, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	configContent := xmlHeader + string(updatedContent)

	os.WriteFile(configurationFile, []byte(configContent), 0644)
	fmt.Println("Updated")
}
