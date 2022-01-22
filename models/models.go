package models

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Coordinates struct {
	Lat string `json:lat`
	Lon string `json:lon`
}

type Weather struct {
	Temperature float32 `json:temperature`
	Pressure    float32 `json:pressure`
	Humidity    float32 `json:humidity`
	Clouds      float32 `json:clouds`
	CurrentTime int64   `json:current_time`
	Sunrise     int64   `json:sunrise`
	Sunset      int64   `json:sunset`
	TimeOffset  int64   `json:timezone_offset`
}

type WeatherTemplateData struct {
	Name        string
	Overall     string
	Temperature int64
	Time        string
	Pressure    float32
	Humidity    float32
	Sunrise     string
	Sunset      string
}

type OopsTemplateData struct {
	RedirectURL string
}

type mongoDBCoordinates struct {
	Name string `json:name`
	Lat  string `json:lat`
	Lon  string `json:lon`
}

func GetCoordinates(searchText string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Print("mongo.Connect() ERROR:", err)
		return "", "", err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("weather").Collection("cities")
	name := strings.ToLower(strings.Replace(searchText, "%20", "", -1))
	filter := bson.D{{"name", name}}
	result := mongoDBCoordinates{}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", "", nil
	}
	return result.Lon, result.Lat, nil
}

func SaveCoordinates(searchText, lat, lon string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("weather").Collection("cities")
	name := strings.ToLower(strings.Replace(searchText, "%20", "", -1))
	toInsert := mongoDBCoordinates{Name: name, Lat: lat, Lon: lon}
	_, err = collection.InsertOne(ctx, toInsert)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}
