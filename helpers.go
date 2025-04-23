package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type responseSuccess struct {
	Body string `json:"body"`
}

type responseError struct {
	Error string `json:"error"`
}

// Simple success response to wrap the correct response body.
func respondSuccess(msg string, w http.ResponseWriter, status int) {
	res := responseSuccess{
		Body: msg,
	}

	writeResponse(res, w, status)
}

// Simple error response to wrap the correct response body.
func respondError(msg string, w http.ResponseWriter, status int) {
	res := responseError{
		Error: msg,
	}

	writeResponse(res, w, status)
}

// DRYs up the handler code, taking any marshallable struct,
// such as Language{} or Word{}, and writing the correct response.
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

// Returns the correct status code, depending on if the failed creation
// was due to a unique constraint violation, or some other unanticipated
// issue.
func getFailedCreationCode(err error) int {
	if strings.HasPrefix(err.Error(), "pq: duplicate key value") {
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}

// Constructs an authenticated endpoint
func (cfg *apiConfig) getAuthenticatedHandler(
	handlerFunc func(w http.ResponseWriter, r *http.Request),
) http.Handler {
	return cfg.auth.AuthenticateAPIKey(
		http.HandlerFunc(handlerFunc),
	)
}
