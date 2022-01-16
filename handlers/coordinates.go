package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"weather_checker/models"

	"context"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// "/coordinates/:cityname"
func CoordinatesHandler(ctx *fasthttp.RequestCtx) {

	lon, lat, status, err := getCoordinates(fmt.Sprintf("%s", ctx.UserValue("cityname")))
	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}
	if status >= 500 {
		ctx.Response.SetStatusCode(500)
		return
	}
	if lat == "" && lon == "" {
		ctx.Response.SetStatusCode(500)
		return
	}
	coords := models.Coordinates{}
	coords.Lat = lat
	coords.Lon = lon

	resp, err := json.Marshal(coords)
	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	ctx.Response.SetStatusCode(200)
	ctx.SetContentType("application/json")
	ctx.Response.SetBody(resp)
}

func getCoordinates(searchText string) (lat, lon string, status int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("weather").Collection("cities")
	name := strings.ToLower(strings.Replace(searchText, "%20", "", -1))
	filter := bson.D{{"name", name}}
	result := models.MongoDBCoordinates{}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == nil {
		return result.Lon, result.Lat, 200, nil
	}

	var geocodingRequest []byte
	URI := geocodeURL + searchText
	status, geocodingRequest, err = fasthttp.Get(geocodingRequest, URI)
	if err != nil {
		log.Printf("error while requesting coordinates from geocoding")
		return "", "", status, err
	}
	if status >= 500 {
		log.Printf("geocode service unvailable")
		return "", "", status, nil
	}
	if status != 200 {
		log.Printf("status not OK in geocoding response")
		return "", "", status, nil
	}

	unmarshaledMap := make(map[string]json.RawMessage)

	err = json.Unmarshal(geocodingRequest, &unmarshaledMap)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return "", "", status, err
	}

	longt := string(unmarshaledMap["longt"][1:])
	latt := string(unmarshaledMap["latt"][1:])
	longt = longt[:len(longt)-1]
	latt = latt[:len(latt)-1]

	toInsert := models.MongoDBCoordinates{Name: name, Lat: latt, Lon: longt}
	_, err = collection.InsertOne(ctx, toInsert)

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
