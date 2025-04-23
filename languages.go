package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"vastestsea/internal/database"

	"github.com/google/uuid"
)

/*
 * Language Handlers
 */

// Get all languages
func (cfg *apiConfig) getLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := cfg.queries.GetLanguages(r.Context())
	if err != nil {
		respondError("No languages found", w, http.StatusNotFound)
		return
	}

	marshallableLanguages := []Language{}
	for _, language := range languages {
		marshallableLanguages = append(marshallableLanguages, getMarshallableLanguage(language))
	}

	writeResponse(marshallableLanguages, w, http.StatusOK)
}

// Get the language specified in the path parameter
func (cfg *apiConfig) getLanguage(w http.ResponseWriter, r *http.Request) {
	languageName := r.PathValue("language")
	language, err := cfg.queries.GetLanguage(r.Context(), strings.ToLower(languageName))
	if err != nil {
		respondError("Language not found", w, http.StatusNotFound)
		return
	}

	writeResponse(getMarshallableLanguage(language), w, http.StatusOK)
}

// Create a new language
// TODO: gate behind authentication
func (cfg *apiConfig) createLanguage(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Name string `json:"name"`
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondError(fmt.Sprintf("Could not decode request body: %s", err), w, http.StatusBadRequest)
		return
	}

	if params.Name == "" {
		respondError("Invalid request body", w, http.StatusBadRequest)
		return
	}

	language, err := cfg.queries.CreateLanguage(r.Context(), params.Name)
	if err != nil {
		respondError(
			fmt.Sprintf("Failed to create language: %s", err),
			w,
			getFailedCreationCode(err),
		)
		return
	}

	writeResponse(getMarshallableLanguage(language), w, http.StatusCreated)
}

func (cfg *apiConfig) updateLanguage(w http.ResponseWriter, r *http.Request) {
	languageName := r.PathValue("language")

	type reqParams struct {
		Name string `json:"name"`
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondError(
			fmt.Sprintf("Could not decode request body: %s", err),
			w,
			http.StatusBadRequest,
		)
		return
	}

	language, err := cfg.queries.UpdateLanguageName(r.Context(), database.UpdateLanguageNameParams{
		Name:   params.Name,
		Name_2: languageName,
	})
	if err != nil {
		language, err = cfg.queries.CreateLanguage(r.Context(), params.Name)
		if err != nil {
			respondError(
				fmt.Sprintf("Failed to create language: %s", err),
				w,
				getFailedCreationCode(err),
			)
			return
		}
		writeResponse(getMarshallableLanguage(language), w, http.StatusCreated)
		return
	}

	writeResponse(getMarshallableLanguage(language), w, http.StatusOK)
}

func (cfg *apiConfig) deleteLanguage(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		ID uuid.UUID `json:"id"`
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondError(
			fmt.Sprintf("Could not decode request body: %s", err),
			w,
			http.StatusBadRequest,
		)
	}

	err := cfg.queries.DeleteLanguage(r.Context(), params.ID)
	if err != nil {
		respondError(
			fmt.Sprintf("Could not delete language: %s", err),
			w,
			http.StatusNotFound,
		)
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
 * Word Handlers
 */

// Get all words registered with a given language, as given in the path parameter
func (cfg *apiConfig) getWordsFromLanguage(w http.ResponseWriter, r *http.Request) {
	languageName := r.PathValue("language")
	language, err := cfg.queries.GetLanguage(r.Context(), strings.ToLower(languageName))
	if err != nil {
		respondError("Language not found", w, http.StatusNotFound)
		return
	}

	words, err := cfg.queries.GetWordsByLanguageID(r.Context(), language.ID)
	if err != nil {
		respondError("No words found", w, http.StatusNotFound)
		return
	}

	marshallableWords := []Word{}
	for _, word := range words {
		marshallableWords = append(marshallableWords, getMarshallableWord(word, []database.Definition{}))
	}

	writeResponse(marshallableWords, w, http.StatusOK)
}

// Get a specific word, as registered in a specific language.
// Both the word and language should be provided in the path parameters.
func (cfg *apiConfig) getWordFromLanguage(w http.ResponseWriter, r *http.Request) {
	languageName := r.PathValue("language")
	language, err := cfg.queries.GetLanguage(r.Context(), strings.ToLower(languageName))
	if err != nil {
		respondError("Language not found", w, http.StatusNotFound)
		return
	}

	wordName := r.PathValue("word")
	word, err := cfg.queries.GetWordFromLanguage(r.Context(), database.GetWordFromLanguageParams{
		Word:       strings.ToLower(wordName),
		LanguageID: language.ID,
	})
	if err != nil {
		respondError("Word not found for that language", w, http.StatusNotFound)
		return
	}

	definitions, _ := cfg.queries.GetDefinitionsOfWord(r.Context(), word.ID)

	writeResponse(getMarshallableWord(word, definitions), w, http.StatusOK)
}

// Get all words registered to any language.
func (cfg *apiConfig) getWords(w http.ResponseWriter, r *http.Request) {
	words, err := cfg.queries.GetWords(r.Context())
	if err != nil {
		log.Println(err.Error())
		respondError("No words found", w, http.StatusNotFound)
		return
	}

	marshallableWords := []Word{}
	for _, word := range words {
		marshallableWord := getMarshallableWord(word, []database.Definition{})

		if definitions, err := cfg.queries.GetDefinitionsOfWord(r.Context(), word.ID); err != nil {
			marshallableWord.Definitions = getMarshallableDefinitions(definitions)
		}

		marshallableWords = append(marshallableWords, marshallableWord)
	}

	writeResponse(marshallableWords, w, http.StatusOK)
}

// Get all possible values of a given word.
// Though words are unique per language, two languages can have words that are
// identical in name. This will provide all variations across all languages,
// should they exist.
func (cfg *apiConfig) getWord(w http.ResponseWriter, r *http.Request) {
	wordName := r.PathValue("word")

	words, err := cfg.queries.GetWord(r.Context(), strings.ToLower(wordName))
	if err != nil {
		respondError("No word found", w, http.StatusNotFound)
		return
	}

	marshallableWords := []Word{}
	for _, word := range words {
		definitions, _ := cfg.queries.GetDefinitionsOfWord(r.Context(), word.ID)

		marshallableWords = append(marshallableWords, getMarshallableWord(word, definitions))
	}

	writeResponse(marshallableWords, w, http.StatusOK)
}

// Create a new word.
// The word itself, and the language of origin, should be provided in the
// request body. If `.language` does not exist in the database, this handler
// creates it, and then creates the word.
func (cfg *apiConfig) createWord(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Word     string `json:"word"`
		Language string `json:"language"`
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondError("Could not decode request body", w, http.StatusBadRequest)
		return
	}

	if params.Word == "" || params.Language == "" {
		respondError("Invalid request body", w, http.StatusBadRequest)
		return
	}

	language, err := cfg.queries.GetLanguage(r.Context(), strings.ToLower(params.Language))
	if err != nil {
		language, err = cfg.queries.CreateLanguage(r.Context(), params.Language)
		if err != nil {
			respondError(
				fmt.Sprintf("Failed to create language: %s", err),
				w,
				getFailedCreationCode(err),
			)
			return
		}
	}

	word, err := cfg.queries.CreateWord(r.Context(), database.CreateWordParams{
		Word:       params.Word,
		LanguageID: language.ID,
	})
	if err != nil {
		respondError(
			fmt.Sprintf("Failed to create word: %s", err),
			w,
			getFailedCreationCode(err),
		)
		return
	}

	writeResponse(getMarshallableWord(word, []database.Definition{}), w, http.StatusCreated)
}

// Create a word for a given language.
// The language should be a path parameter. The word should be provided in the body.
func (cfg *apiConfig) createWordForLanguage(w http.ResponseWriter, r *http.Request) {
	languageName := r.PathValue("language")
	language, err := cfg.queries.GetLanguage(r.Context(), languageName)
	if err != nil {
		respondError("Language not found", w, http.StatusNotFound)
		return
	}

	type reqParams struct {
		Word string `json:"word"`
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondError("Could not decode request body", w, http.StatusInternalServerError)
		return
	}

	if params.Word == "" {
		respondError("Invalid request body", w, http.StatusBadRequest)
		return
	}

	word, err := cfg.queries.CreateWord(r.Context(), database.CreateWordParams{
		Word:       params.Word,
		LanguageID: language.ID,
	})
	if err != nil {
		respondError(
			fmt.Sprintf("Failed to create word: %s", err),
			w,
			getFailedCreationCode(err),
		)
		return
	}

	writeResponse(getMarshallableWord(word, []database.Definition{}), w, http.StatusCreated)
}
