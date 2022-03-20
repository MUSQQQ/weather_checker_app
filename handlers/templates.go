package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"weather_checker/models"

	"github.com/valyala/fasthttp"
)

// "/weather/:cityname"
func WeatherPageHandler(ctx *fasthttp.RequestCtx) {

	lon, lat, status, err := getCoordinates(fmt.Sprintf("/weather/%s", ctx.UserValue("cityname")))
	if err != nil || lat == "" || lon == "" {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if status != http.StatusOK {
		oopsData := models.OopsTemplateData{
			RedirectURL: fmt.Sprintf("/weather/%s", ctx.UserValue("cityname")),
		}
		tpl := template.Must(template.ParseFiles("templates/oops.gohtml"))
		ctx.SetContentType("text/html")
		tpl.Execute(ctx, oopsData)
		return
	}

	weather, status, err := getWeatherData(lat, lon)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if status != http.StatusOK {
		oopsData := models.OopsTemplateData{
			RedirectURL: fmt.Sprintf("/weather/%s", ctx.UserValue("cityname")),
		}
		tpl := template.Must(template.ParseFiles("templates/oops.gohtml"))
		ctx.SetContentType("text/html")
		tpl.Execute(ctx, oopsData)
		return
	}

	name := strings.ToUpper(strings.Replace(fmt.Sprintf("%s", ctx.UserValue("cityname")), "%20", " ", -1))

	templateData := models.WeatherTemplateData{
		Name:        name,
		Overall:     determineOverall(weather),
		Temperature: int64(weather.Temperature),
		Time:        timestampToString(time.Unix(weather.CurrentTime, 0).Unix(), weather.TimeOffset),
		Pressure:    weather.Pressure,
		Humidity:    weather.Humidity,
		Sunrise:     strings.Split(timestampToString(time.Unix(weather.Sunrise, 0).Unix(), weather.TimeOffset), "\n")[1],
		Sunset:      strings.Split(timestampToString(time.Unix(weather.Sunset, 0).Unix(), weather.TimeOffset), "\n")[1],
	}

	tpl := template.Must(template.ParseFiles("templates/weather.gohtml"))
	ctx.SetContentType("text/html")
	tpl.Execute(ctx, templateData)
}
