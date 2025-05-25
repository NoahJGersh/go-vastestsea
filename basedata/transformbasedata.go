package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
)

func main() {
	rawData, err := os.ReadFile("./sources/baselexdata.json")
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	type wordData struct {
		NameRoman  string `json:"name_roman"`
		NameGlyph  string `json:"name_glyph"`
		Part       string `json:"part"`
		Definition string `json:"definition"`
	}

	type lexData struct {
		Glyph string `json:"glyph"`
		Words []wordData
	}

	type lexicon map[string]lexData

	baseLexicon := make(lexicon)
	if err = json.Unmarshal(rawData, &baseLexicon); err != nil {
		log.Fatalf("Failed to unmarshal data: %s", err)
	}

	type finalWordParamsDef struct {
		Content      string `json:"content"`
		PartOfSpeech string `json:"part_of_speech"`
	}

	type finalWordParams struct {
		Word       string             `json:"word"`
		Formatted  string             `json:"formatted"`
		Definition finalWordParamsDef `json:"definition"`
	}

	finalLexData := make([]finalWordParams, 0)

	for _, lexSection := range baseLexicon {
		for _, word := range lexSection.Words {
			wordParams := finalWordParams{
				Word:      word.NameRoman,
				Formatted: word.NameGlyph,
				Definition: finalWordParamsDef{
					Content:      word.Definition,
					PartOfSpeech: word.Part,
				},
			}

			finalLexData = append(finalLexData, wordParams)
		}
	}

	// Write final data to disk
	finalData, err := json.Marshal(finalLexData)
	if err != nil {
		log.Fatalf("Unable to marshal final data: %s", err)
	}

	err = os.WriteFile("./sources/finallexdata.json", finalData, fs.ModePerm)
	if err != nil {
		log.Fatalf("Unable to write data to file: %s", err)
	}

	log.Println("Data transformed and written to disk")
}
