package handlers

import (
	"encoding/json"
	"log"

	"weather_checker/models"

	"github.com/valyala/fasthttp"
)

func getCoordinates(searchText string) (longt, latt string, status int, err error) {
	lat, lon, err := models.GetCoordinates(searchText)
	if err == nil && lat != "" && lon != "" {
		return lat, lon, 200, nil
	}
	err = nil

	var geocodingRequest []byte
	URI := geocodeURL + searchText
	status, geocodingRequest, err = fasthttp.Get(geocodingRequest, URI)
	if err != nil {
		log.Print(err)
		log.Println(string(geocodingRequest))
		return "", "", status, err
	}
	if status != 200 {
		log.Printf("status not OK in geocoding response:")
		log.Println(string(geocodingRequest))
		return "", "", status, nil
	}

	unmarshaledMap := make(map[string]json.RawMessage)

	err = json.Unmarshal(geocodingRequest, &unmarshaledMap)
	if err != nil {
		log.Print(err)
		return "", "", status, err
	}

	longt = string(unmarshaledMap["longt"][1:])
	latt = string(unmarshaledMap["latt"][1:])
	longt = longt[:len(longt)-1]
	latt = latt[:len(latt)-1]

	models.SaveCoordinates(searchText, latt, longt)

	return longt, latt, 200, nil
}
