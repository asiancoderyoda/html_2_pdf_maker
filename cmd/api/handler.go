package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) serverStatus(w http.ResponseWriter, r *http.Request) {
	status := ServerStatus{
		Version: GetEnvFromKey("VERSION"),
		Status:  "OK",
		Env:     app.config.env,
	}

	err := app.writeJSON(w, http.StatusOK, status, "server_status")
	if err != nil {
		log.Fatalf("Error while writing response: %v", err)
		return
	}
}

/*
TODO:
Add functionality to upload generated pdf to S3 bucket
*/

/*
TODO:
Add functionality to enqueue a audit event to sqs
*/
func (app *Application) createInvoice(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	templateType := params.ByName("type")
	save2bucket := params.ByName("save2bucket")

	fmt.Println("templateType: ", templateType, " save2bucket: ", save2bucket)

	request := Request{}
	json.NewDecoder(r.Body).Decode(&request)

	parsedTemplateData, err := app.fetchTemplate(templateType, request.Data)

	if err != nil {
		app.writeError(w, err)
		return
	}

	pathToFile, err := app.parseTemplate(templateType, parsedTemplateData)

	if err != nil {
		app.writeError(w, err)
		return
	}

	fmt.Println("Generating PDF for file: ", pathToFile)

	generatedPdfPath, err := app.wkhtmltopdf.createPdf(pathToFile)

	if err != nil {
		app.writeError(w, err)
		return
	}

	fileBytes, err := ioutil.ReadFile(generatedPdfPath)

	if err != nil {
		app.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=invoice.pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	w.Write(fileBytes)
}
