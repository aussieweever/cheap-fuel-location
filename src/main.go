package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"njetch.com/cheap-fuel-location.git/src/models"
	xmlmodel "njetch.com/cheap-fuel-location.git/src/models/xml"
)

var xmlHeader = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	err = godotenv.Load(exPath + "/.env")
	if err != nil {
		// mainly for development purposes
		err = godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	fuelType := os.Getenv("FUEL_TYPE")
	fmt.Printf("Fuel type: %s\n", fuelType)
	state := os.Getenv("STATE")

	fmt.Print("State: " + state + "\n")

	configurationFile := os.Getenv("FILE_PATH")
	fmt.Println("Configuration file: " + configurationFile)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("redirect to: ", req.URL)
			return nil
		},
	}

	req, err := http.NewRequest("GET", "https://projectzerothree.info/api.php?format=json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error to request price list", err)
		return
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error to read price list", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		redirectUrl, err := resp.Location()
		if err != nil {
			fmt.Println("Error getting redirect location", err)
			return
		}

		req.URL = redirectUrl
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Error to request price list from the redirect url", redirectUrl, err)
			return
		}
		content, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error to read price list", err)
			return
		}
		defer resp.Body.Close()
	} else if resp.StatusCode >= 400 {
		fmt.Println("Error to request price list due to error status code: ", resp.StatusCode)
		return
	}

	var queryResult models.QueryResult

	json.Unmarshal(content, &queryResult)
	var priceItems []models.PriceItem

	for _, region := range queryResult.Regions {
		if state == "" || state == region.Region {
			priceItems = append(priceItems, region.Prices...)
		}
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

	if cheapestPriceItem == nil {
		fmt.Println("No prices found")
		return
	}

	fmt.Printf("The cheapest diesel is from %s %s at %f\n", cheapestPriceItem.Name, cheapestPriceItem.State, cheapestPriceItem.Price)
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
