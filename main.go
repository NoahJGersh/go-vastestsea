package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"vastestsea/internal/auth"
	"vastestsea/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	queries *database.Queries
	auth    auth.AuthConfig
}

func main() {
	// Setup db connection
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to establish connection to database. Exiting.")
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		queries: dbQueries,
		auth: auth.AuthConfig{
			ApiKey: os.Getenv("API_KEY"),
		},
	}

	// Construct mux
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /vs/languages", apiCfg.getLanguages)
	serveMux.HandleFunc("GET /vs/languages/{language}", apiCfg.getLanguage)
	serveMux.HandleFunc("GET /vs/languages/{language}/words", apiCfg.getWordsFromLanguage)
	serveMux.HandleFunc("GET /vs/languages/{language}/words/{word}", apiCfg.getWordFromLanguage)
	serveMux.HandleFunc("GET /vs/languages/words", apiCfg.getWords)
	serveMux.HandleFunc("GET localhost/vs/languages/words/{word}", apiCfg.getWord)

	// Authenticated endpoints
	serveMux.Handle("POST /vs/languages", apiCfg.getAuthenticatedHandler(apiCfg.createLanguage))
	serveMux.Handle("DELETE /vs/languages", apiCfg.getAuthenticatedHandler(apiCfg.deleteLanguage))
	serveMux.Handle("PUT /vs/languages/{language}", apiCfg.getAuthenticatedHandler(apiCfg.updateLanguage))
	serveMux.Handle("POST /vs/languages/{language}/words", apiCfg.getAuthenticatedHandler(apiCfg.createWordForLanguage))
	serveMux.Handle("POST /vs/languages/words", apiCfg.getAuthenticatedHandler(apiCfg.createWord))

	// Run server
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	log.Fatal(server.ListenAndServe())
}
