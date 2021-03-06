package main

import (
	"html/template"
	"log"
	"net/http"

	"weather_checker/handlers"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {

	go func() {
		http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
		http.HandleFunc("/", Index)
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	r := router.New()
	r.GET("/", handlers.MainPageHandler)
	r.GET("/about", handlers.AboutPageHandler)
	r.GET("/weather/{cityname}", handlers.WeatherPageHandler)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, nil)
}
