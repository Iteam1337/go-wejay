package jsonResponses

import (
	"encoding/json"
	"errors"
	"net/http"
)

// PlayerResponse …
type PlayerResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func (p *PlayerResponse) Write(w http.ResponseWriter) {
	json, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// NewPlayerResponse …
func NewPlayerResponse(w http.ResponseWriter, message string) (res PlayerResponse) {
	res.Ok = true
	res.Message = message
	res.Write(w)
	return
}

// NewPlayerResponseErr …
func NewPlayerResponseErr(w http.ResponseWriter, message string) (res PlayerResponse) {
	res.Ok = false
	res.Message = message
	res.Error = errors.New(message)
	res.Write(w)
	return
}
