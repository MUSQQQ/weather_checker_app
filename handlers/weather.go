package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"weather_checker/models"

	"github.com/valyala/fasthttp"
)

func getWeatherData(lat string, lon string) (weather *models.Weather, status int, err error) {
	var openWeatherRequest []byte
	URI := fmt.Sprintf("%s?lat=%s&lon=%s&exclude=%s&appid=%s", openWeatherURL, lat, lon, toExclude, openWeatherAPIKey)

	status, openWeatherRequest, err = fasthttp.Get(openWeatherRequest, URI)
	if err != nil {
		log.Printf("error while requesting coordinates from geocoding")
		return nil, 500, err
	}
	if status != 200 {
		log.Printf("openweather service unvailable or wrong request")
		log.Print(string(openWeatherRequest))
		return nil, status, nil
	}
	unmarshaledMap1 := make(map[string]json.RawMessage)

	err = json.Unmarshal(openWeatherRequest, &unmarshaledMap1)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return nil, 590, err
	}
	unmarshaledMap2 := make(map[string]json.RawMessage)

	err = json.Unmarshal(unmarshaledMap1["current"], &unmarshaledMap2)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return nil, 590, err
	}
	weather = &models.Weather{}
	err = json.Unmarshal(unmarshaledMap1["timezone_offset"], &weather.TimeOffset)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return nil, 590, err
	}

	weather.Temperature, err = byteArrayToFloat(unmarshaledMap2["temp"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return nil, status, err
	}
	weather.Temperature -= 272.15

	weather.Pressure, err = byteArrayToFloat(unmarshaledMap2["pressure"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return nil, status, err
	}
	weather.Humidity, err = byteArrayToFloat(unmarshaledMap2["humidity"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return nil, status, err
	}
	weather.Clouds, err = byteArrayToFloat(unmarshaledMap2["clouds"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return nil, status, err
	}
	err = json.Unmarshal(unmarshaledMap2["dt"], &weather.CurrentTime)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return nil, 590, err
	}
	err = json.Unmarshal(unmarshaledMap2["sunrise"], &weather.Sunrise)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return nil, 590, err
	}
	err = json.Unmarshal(unmarshaledMap2["sunset"], &weather.Sunset)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return nil, 590, err
	}

	return weather, status, nil
}
