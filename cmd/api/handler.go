package main

import (
	"encoding/json"
	"fmt"
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
	var invoiceData Invoice
	err := json.NewDecoder(r.Body).Decode(&invoiceData)
	if err != nil {
		app.writeError(w, err)
		return
	}

	pathToFile, err := app.parseTemplate(invoiceData)

	if err != nil {
		app.writeError(w, err)
		return
	}

	fmt.Println("Generating PDF for file: ", pathToFile)

	pdf, err := app.wkhtmltopdf.createPdf(pathToFile)

	if err != nil {
		app.writeError(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, pdf, "pdf")

	if err != nil {
		log.Fatalf("Error while writing response: %v", err)
		return
	}
}
