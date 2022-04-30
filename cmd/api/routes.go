package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("GET", "/server/status", app.serverStatus)
	router.HandlerFunc("POST", "/api/v1/:type/:save2bucket", app.html2pdf)
	router.HandlerFunc("GET", "/api/v1/:type/:key", app.fetchS3Item)
	return app.enableCors(router)
}
