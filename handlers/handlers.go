package handlers

import (
	"html/template"

	"github.com/valyala/fasthttp"
)

type aboutData struct {
	Back string
}

func MainPageHandler(ctx *fasthttp.RequestCtx) {
	tpl := template.Must(template.ParseFiles("templates/index.gohtml"))
	ctx.SetContentType("text/html")
	tpl.Execute(ctx, nil)
}

func AboutPageHandler(ctx *fasthttp.RequestCtx) {
	tpl := template.Must(template.ParseFiles("templates/about.gohtml"))
	ctx.SetContentType("text/html")
	ad := aboutData{"Home"}
	tpl.Execute(ctx, ad)
}
