package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type responseSuccess struct {
	Body string `json:"body"`
}

type responseError struct {
	Error string `json:"error"`
}

func respondSuccess(msg string, w http.ResponseWriter, status int) {
	res := responseSuccess{
		Body: msg,
	}

	writeResponse(res, w, status)
}

func respondError(msg string, w http.ResponseWriter, status int) {
	res := responseError{
		Error: msg,
	}

	writeResponse(res, w, status)
}

func writeResponse[T any](res T, w http.ResponseWriter, status int) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
