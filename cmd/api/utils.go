package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

func (app *Application) writeJSON(w http.ResponseWriter, statusCode int, payload interface{}, wrap string) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = payload
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)

	return nil
}

func (app *Application) writeError(w http.ResponseWriter, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	httpError := errorResponse{
		Error: err.Error(),
	}

	app.writeJSON(w, http.StatusUnprocessableEntity, httpError, "error")

}

func (app *Application) parseTemplate(data Invoice) (string, error) {
	var document bytes.Buffer // buffer to hold the final document

	// Load the HTML template
	templatePath := fmt.Sprintf("%s%s%s", app.config.templateDir, "invoice", app.config.htmlExtension)
	tmpl, err := template.ParseFiles(templatePath)

	if err != nil {
		return "", err
	}

	// Execute the template
	err = tmpl.Execute(&document, data)

	if err != nil {
		return "", err
	}

	// Create populated HTML template
	populatedTemplate := fmt.Sprintf("%s%d-%d%s", app.config.tempDir, data.Id, int32(time.Now().UnixNano()), app.config.htmlExtension)
	file, _ := os.Create(populatedTemplate)
	defer file.Close()

	// Write the populated HTML template to file
	file.Write(document.Bytes())

	return populatedTemplate, nil
}