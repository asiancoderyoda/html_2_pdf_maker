package main

import (
	"log"
	"net/http"
)

func (app *Application) serverStatus(w http.ResponseWriter, r *http.Request) {
	status := ServerStatus{
		Version: version,
		Status:  "OK",
		Env:     app.config.env,
	}

	err := app.writeJSON(w, http.StatusOK, status, "server_status")
	if err != nil {
		log.Fatalf("Error while writing response: %v", err)
		return
	}
}

func (app *Application) createInvoice(w http.ResponseWriter, r *http.Request) {

}
