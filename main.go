package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"vastestsea/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	queries *database.Queries
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
	}

	// Construct mux
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /vs/languages", apiCfg.getLanguages)
	serveMux.HandleFunc("POST /vs/languages", apiCfg.createLanguage)
	serveMux.HandleFunc("GET /vs/languages/{language}", apiCfg.getLanguage)
	serveMux.HandleFunc("GET /vs/languages/{language}/words", apiCfg.getWordsFromLanguage)
	serveMux.HandleFunc("POST /vs/languages/{language}/words", apiCfg.createWordForLanguage)
	serveMux.HandleFunc("GET /vs/languages/{language}/words/{word}", apiCfg.getWordFromLanguage)
	serveMux.HandleFunc("GET /vs/languages/words", apiCfg.getWords)
	serveMux.HandleFunc("POST /vs/languages/words", apiCfg.createWord)
	serveMux.HandleFunc("GET localhost/vs/languages/words/{word}", apiCfg.getWord)

	// Run server
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	log.Fatal(server.ListenAndServe())
}
