package main

import (
	"time"
	"vastestsea/internal/database"

	"github.com/google/uuid"
)

type Language struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getMarshallableLanguage(l database.Language) Language {
	marshallable := Language{
		ID:        l.ID,
		Name:      l.Name,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}

	return marshallable
}

type Word struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Word          string    `json:"word"`
	FontFormatted string    `json:"font_formatted"`
	LanguageID    uuid.UUID `json:"language_id"`
}

func getMarshallableWord(w database.Word) Word {
	marshallable := Word{
		ID:         w.ID,
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
		Word:       w.Word,
		LanguageID: w.LanguageID,
	}

	if w.FontFormatted.Valid {
		marshallable.FontFormatted = w.FontFormatted.String
	}

	return marshallable
}
