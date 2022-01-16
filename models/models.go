package models

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

type MongoDBCoordinates struct {
	Name string `json:name`
	Lat  string `json:lat`
	Lon  string `json:lon`
}
