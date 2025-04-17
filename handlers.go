package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
 * Languages
 */

func (cfg *apiConfig) getLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := cfg.queries.GetLanguages(r.Context())
	if err != nil {
		respondError("No languages found", w, http.StatusNotFound)
		return
	}

	marshallableLanguages := []Language{}
	for _, language := range languages {
		marshallableLanguages = append(marshallableLanguages, Language{
			ID:        language.ID,
			Name:      language.Name,
			CreatedAt: language.CreatedAt,
			UpdatedAt: language.UpdatedAt,
		})
	}

	writeResponse(marshallableLanguages, w, http.StatusOK)
}

func (cfg *apiConfig) getLanguage(w http.ResponseWriter, r *http.Request) {
	languageName := r.PathValue("language")
	language, err := cfg.queries.GetLanguage(r.Context(), strings.ToLower(languageName))
	if err != nil {
		respondError("Language not found", w, http.StatusNotFound)
		return
	}

	marshallableLanguage := Language{
		ID:        language.ID,
		Name:      language.Name,
		CreatedAt: language.CreatedAt,
		UpdatedAt: language.UpdatedAt,
	}

	writeResponse(marshallableLanguage, w, http.StatusOK)
}

func (cfg *apiConfig) createLanguage(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name string `json:"name"`
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondError("Could not decode request body", w, http.StatusBadRequest)
	}

	language, err := cfg.queries.CreateLanguage(r.Context(), params.Name)
	if err != nil {
		respondError(
			fmt.Sprintf("Failed to create language: %s", err),
			w,
			http.StatusInternalServerError,
		)
		return
	}

	marshallableLanguage := Language{
		ID:        language.ID,
		Name:      language.Name,
		CreatedAt: language.CreatedAt,
		UpdatedAt: language.UpdatedAt,
	}

	writeResponse(marshallableLanguage, w, http.StatusCreated)
}
