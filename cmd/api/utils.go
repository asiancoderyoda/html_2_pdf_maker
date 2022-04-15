package main

import (
	"encoding/json"
	"net/http"
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
