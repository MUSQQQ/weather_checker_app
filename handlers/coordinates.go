package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"weather_checker/models"

	"github.com/valyala/fasthttp"
)

// "/coordinates/:cityname"
func CoordinatesHandler(ctx *fasthttp.RequestCtx) {

	lon, lat, status, err := getCoordinates(fmt.Sprintf("%s", ctx.UserValue("cityname")))
	if err != nil || status >= 500 || (lat == "" && lon == "") {
		ctx.Response.SetStatusCode(500)
		return
	}
	coords := models.Coordinates{
		Lat: lat,
		Lon: lon,
	}

	resp, err := json.Marshal(coords)
	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	ctx.Response.SetStatusCode(200)
	ctx.SetContentType("application/json")
	ctx.Response.SetBody(resp)
}

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
		log.Printf("error while requesting coordinates from geocoding")
		log.Println(string(geocodingRequest))
		return "", "", status, err
	}
	if status >= 500 {
		log.Printf("geocode service unvailable")
		log.Println(string(geocodingRequest))
		return "", "", status, nil
	}
	if status != 200 {
		log.Printf("status not OK in geocoding response:")
		log.Println(string(geocodingRequest))
		return "", "", status, nil
	}

	unmarshaledMap := make(map[string]json.RawMessage)

	err = json.Unmarshal(geocodingRequest, &unmarshaledMap)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return "", "", status, err
	}

	longt = string(unmarshaledMap["longt"][1:])
	latt = string(unmarshaledMap["latt"][1:])
	longt = longt[:len(longt)-1]
	latt = latt[:len(latt)-1]

	models.SaveCoordinates(searchText, latt, longt)

	return longt, latt, 200, nil
}

func byteArrayToFloat(bytes []byte) (result float32, err error) {
	strByte := string(bytes)
	var i int
	for i = 0; i < len(strByte); i++ {
		if strByte[i] == '.' {
			break
		}
	}

	var intResultPart int

	intResultPart, err = strconv.Atoi(strByte[0:i])

	if err != nil {
		log.Printf("error while converting: %s", err)
		return -1, err
	}

	var mantissaPart int
	if i >= len(strByte)-1 {
		mantissaPart = 0

	} else {
		mantissaPart, err = strconv.Atoi(strByte[i+1:])
	}

	if err != nil {
		log.Printf("error while converting: %s", err)
		return -1, err
	}
	result = float32(intResultPart) + float32(mantissaPart)/(float32(math.Pow10(len(strByte)-i-1)))
	return result, nil
}
