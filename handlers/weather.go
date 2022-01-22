package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"weather_checker/models"

	"github.com/valyala/fasthttp"
)

// "/weather/checkcity/:cityname"
func MainWeatherHandler(ctx *fasthttp.RequestCtx) {

	//https://api.weather.gov/points/

	longt, latt, status, err := getCoordinates(fmt.Sprintf("%s", ctx.UserValue("cityname")))
	if err != nil || status >= 500 || latt == "" && longt == "" {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	temp, pressure, humidity, clouds, currentTime, sunrise, sunset, offset, status, err := getWeatherData(latt, longt)
	if err != nil || status != 200 {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	temp -= 272.15 //convert to Celsius

	clientTime := time.Unix(currentTime, 0)
	clientSunrise := time.Unix(sunrise, 0)
	clientSunset := time.Unix(sunset, 0)

	weather := models.Weather{}
	weather.Humidity = humidity
	weather.Pressure = pressure
	weather.Temperature = temp
	weather.Clouds = clouds
	weather.CurrentTime = clientTime.Unix()
	weather.Sunrise = clientSunrise.Unix()
	weather.Sunset = clientSunset.Unix()
	weather.TimeOffset = offset

	resp, err := json.Marshal(weather)
	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}
	ctx.Response.SetStatusCode(200)
	ctx.SetContentType("application/json")
	ctx.Response.SetBody(resp)
}

func getWeatherData(lat string, lon string) (temp, pressure, humidity, clouds float32, currentTime, sunrise, sunset, offset int64, status int, err error) {
	toExclude := "minutely,hourly,daily,alerts"
	var openWeatherRequest []byte
	URI := fmt.Sprintf("%s?lat=%s&lon=%s&exclude=%s&appid=%s", openWeatherURL, lat, lon, toExclude, openWeatherAPIKey)

	status, openWeatherRequest, err = fasthttp.Get(openWeatherRequest, URI)

	if err != nil {
		log.Printf("error while requesting coordinates from geocoding")
		return 0, 0, 0, 0, 0, 0, 0, 0, 500, err
	}
	if status != 200 {
		log.Printf("openweather service unvailable or wrong request")
		log.Print(string(openWeatherRequest))
		return 0, 0, 0, 0, 0, 0, 0, 0, status, nil
	}
	unmarshaledMap1 := make(map[string]json.RawMessage)

	err = json.Unmarshal(openWeatherRequest, &unmarshaledMap1)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 0, 0, 0, 590, err
	}
	unmarshaledMap2 := make(map[string]json.RawMessage)

	err = json.Unmarshal(unmarshaledMap1["current"], &unmarshaledMap2)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 0, 0, 0, 590, err
	}
	err = json.Unmarshal(unmarshaledMap1["timezone_offset"], &offset)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 0, 0, 0, 590, err
	}

	temp, err = byteArrayToFloat(unmarshaledMap2["temp"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return 0, 0, 0, 0, 0, 0, 0, 0, status, err
	}

	pressure, err = byteArrayToFloat(unmarshaledMap2["pressure"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return 0, 0, 0, 0, 0, 0, 0, 0, status, err
	}
	humidity, err = byteArrayToFloat(unmarshaledMap2["humidity"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return 0, 0, 0, 0, 0, 0, 0, 0, status, err
	}
	clouds, err = byteArrayToFloat(unmarshaledMap2["clouds"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return 0, 0, 0, 0, 0, 0, 0, 0, status, err
	}
	err = json.Unmarshal(unmarshaledMap2["dt"], &currentTime)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 0, 0, 0, 590, err
	}
	err = json.Unmarshal(unmarshaledMap2["sunrise"], &sunrise)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 0, 0, 0, 590, err
	}
	err = json.Unmarshal(unmarshaledMap2["sunset"], &sunset)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 0, 0, 0, 590, err
	}

	return temp, pressure, humidity, clouds, currentTime, sunrise, sunset, offset, status, nil
}
