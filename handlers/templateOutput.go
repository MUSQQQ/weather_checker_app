package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"weather_checker/models"

	"github.com/valyala/fasthttp"
)

// /weather/:cityname
func WeatherPageHandler(ctx *fasthttp.RequestCtx) {

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/weather/checkcity/%s", ctx.UserValue("cityname")))
	if err != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != 200 {
		oopsData := models.OopsTemplateData{
			RedirectURL: fmt.Sprintf("/weather/%s", ctx.UserValue("cityname")),
		}
		tpl := template.Must(template.ParseFiles("templates/oops.gohtml"))
		ctx.SetContentType("text/html")
		tpl.Execute(ctx, oopsData)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}
	bodyStr := models.Weather{}
	err = json.Unmarshal(body, &bodyStr)
	if err != nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}

	name := strings.ToUpper(strings.Replace(fmt.Sprintf("%s", ctx.UserValue("cityname")), "%20", " ", -1))

	templateData := models.WeatherTemplateData{
		Name:        name,
		Overall:     determineOverall(bodyStr),
		Temperature: int64(bodyStr.Temperature),
		Time:        timestampToString(bodyStr.CurrentTime, bodyStr.TimeOffset),
		Pressure:    bodyStr.Pressure,
		Humidity:    bodyStr.Humidity,
		Sunrise:     strings.Split(timestampToString(bodyStr.Sunrise, bodyStr.TimeOffset), "\n")[1],
		Sunset:      strings.Split(timestampToString(bodyStr.Sunset, bodyStr.TimeOffset), "\n")[1],
	}
	tpl := template.Must(template.ParseFiles("templates/weather.gohtml"))
	ctx.SetContentType("text/html")
	tpl.Execute(ctx, templateData)
}

func determineOverall(data models.Weather) (overall string) {
	switch {
	case data.Temperature < 0.0:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cold, possible rain"
					default:
						overall = "Cold and sunny"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cold, cloudy and possible rain"
					default:
						overall = "Cold and cloudy"
					}

				}

			}

		}
	case data.Temperature < 10.0:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cool, possible rain"
					default:
						overall = "Cool and sunny"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Cool, cloudy and possible rain"
					default:
						overall = "Cool and cloudy"
					}

				}

			}
		}
	case data.Temperature < 20.0:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Warm, possible rain"
					default:
						overall = "Warm and sunny"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Warm, cloudy and possible rain"
					default:
						overall = "Warm and cloudy"
					}

				}

			}

		}
	default:
		{
			switch {
			case data.Clouds < 50.0:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Hot, possible rain"
					default:
						overall = "Hot and sunny"
					}

				}
			default:
				{
					switch {
					case data.Humidity > 90.0:
						overall = "Hot, cloudy and possible rain"
					default:
						overall = "Hot and cloudy"
					}

				}

			}
		}
	}
	return overall
}

func timestampToString(timestamp, offset int64) (result string) {
	year, month, day := time.Unix(timestamp, 0).UTC().Add(time.Duration(offset) * time.Second).Date()
	hour, minute, _ := time.Unix(timestamp, 0).UTC().Add(time.Duration(offset) * time.Second).Clock()

	result = fmt.Sprintf("%d:%d:%d\n%d:%d", year, month, day, hour, minute)
	return result
}
