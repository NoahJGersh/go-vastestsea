package main

import (
	"net/http"
)

func languageHandler(w http.ResponseWriter, _ *http.Request) {
	respondSuccess("Hello lang!", w, http.StatusOK)
}

func languagesHandler(w http.ResponseWriter, _ *http.Request) {
	respondError("Not yet implemented", w, http.StatusNotImplemented)
}
