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
	ID            uuid.UUID    `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Word          string       `json:"word"`
	FontFormatted string       `json:"font_formatted"`
	LanguageID    uuid.UUID    `json:"language_id"`
	Definitions   []Definition `json:"definitions,omitempty"`
}

func getMarshallableWord(w database.Word, d []database.Definition) Word {
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

	if len(d) > 0 {
		marshallable.Definitions = getMarshallableDefinitions(d)
	}

	return marshallable
}

type Definition struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Content      string    `json:"content"`
	PartOfSpeech string    `json:"part_of_speech"`
	WordID       uuid.UUID `json:"word_id"`
}

func getMarshallableDefinition(d database.Definition) Definition {
	marshallable := Definition{
		ID:           d.ID,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		Content:      d.Content,
		PartOfSpeech: d.PartOfSpeech,
		WordID:       d.WordID,
	}

	return marshallable
}

// Not a strictly necessary function, but it helps to avoid over-nesting when marshalling
// words with definitions
func getMarshallableDefinitions(definitions []database.Definition) []Definition {
	marshallable := []Definition{}

	for _, d := range definitions {
		marshallable = append(marshallable, getMarshallableDefinition(d))
	}

	return marshallable
}
