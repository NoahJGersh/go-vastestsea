package main

import (
	"fmt"
	"log"
	"net/http"
	"vastestsea/internal/auth"
	"vastestsea/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	queries  *database.Queries
	auth     auth.AuthConfig
	hostName string
}

func main() {
	// Initialize environment and db connection
	apiCfg := initializeConfig()

	// Construct mux
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /vs/languages", apiCfg.getLanguages)
	serveMux.HandleFunc("GET /vs/languages/{language}", apiCfg.getLanguage)
	serveMux.HandleFunc("GET /vs/languages/{language}/words", apiCfg.getWordsFromLanguage)
	serveMux.HandleFunc("GET /vs/languages/{language}/words/{word}", apiCfg.getWordFromLanguage)
	serveMux.HandleFunc("GET /vs/languages/words", apiCfg.getWords)
	serveMux.HandleFunc(
		fmt.Sprintf("GET %s/vs/languages/words/{word}", apiCfg.hostName),
		apiCfg.getWord,
	)

	// Authenticated endpoints
	serveMux.Handle("POST /vs/languages", apiCfg.getAuthenticatedHandler(apiCfg.createLanguage))
	serveMux.Handle("DELETE /vs/languages", apiCfg.getAuthenticatedHandler(apiCfg.deleteLanguage))
	serveMux.Handle("PUT /vs/languages/{language}", apiCfg.getAuthenticatedHandler(apiCfg.updateLanguage))
	serveMux.Handle("POST /vs/languages/{language}/words", apiCfg.getAuthenticatedHandler(apiCfg.createWordForLanguage))
	serveMux.Handle("PUT /vs/languages/{language}/words/{word}", apiCfg.getAuthenticatedHandler(apiCfg.updateWord))
	serveMux.Handle("DELETE /vs/languages/{language}/words/{word}", apiCfg.getAuthenticatedHandler(apiCfg.deleteWordFromLanguage))
	serveMux.Handle("POST /vs/languages/words", apiCfg.getAuthenticatedHandler(apiCfg.createWords))
	serveMux.Handle("DELETE /vs/languages/words/{word_id}", apiCfg.getAuthenticatedHandler(apiCfg.deleteWord))

	// Run server
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	log.Fatal(server.ListenAndServe())
}
