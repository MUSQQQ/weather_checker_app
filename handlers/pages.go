package handlers

import (
	"html/template"

	"weather_checker/models"

	"github.com/valyala/fasthttp"
)

func MainPageHandler(ctx *fasthttp.RequestCtx) {
	tpl := template.Must(template.ParseFiles("templates/index.gohtml"))
	ctx.SetContentType("text/html")
	tpl.Execute(ctx, nil)
}

func AboutPageHandler(ctx *fasthttp.RequestCtx) {
	tpl := template.Must(template.ParseFiles("templates/about.gohtml"))
	ctx.SetContentType("text/html")
	ad := models.AboutData{
		Back: "Home",
	}
	tpl.Execute(ctx, ad)
}
