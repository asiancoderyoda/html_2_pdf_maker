package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("GET", "/server/status", app.serverStatus)
	router.HandlerFunc("POST", "/api/v1/invoices", app.createInvoice)
	return app.enableCors(router)
}
