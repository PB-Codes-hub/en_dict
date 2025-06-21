package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Structure matching dictionaryapi.dev response

type DectionaryEntry struct {
	Word string `json:"word"`
	Phonetic string `json:"phonetic"`
	Meanings []struct {
		PartOfSpeech string `json:"partofspeech"`
		Definitions []struct {
			Definition string `json:"definition"`
			Example string `json:"example"`
			Synonyms []string `json:"synonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func main() {
	// Check if the user has provide a word
	// As and argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: en_dict <word>")
		os.Exit(1) // Status 1 Error
	}

	word := os.Args[1]

	// API request URL
	url := "https://api.dictionaryapi.dev/api/v2/entries/en/" + word

	// Make a GET request
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error during HTTP request:", err)
		os.Exit(1)
	}

	// Close Response body on func exit
	defer resp.Body.Close()

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	// If status code is not 200 (OK), no defination found
	if resp.StatusCode != 200 {
		fmt.Printf("No defination found for: %s\n", word)
		os.Exit(1)
	}

	// The API returns a JSON array of definations,
	var results []DectionaryEntry

	// Decode JSON into GO structs
	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// We only care about the first entry,
	// Usually the best match
	entry := results[0]

	// Print the word and its phonetic spelling (if available)
	fmt.Printf("📖 %s", entry.Word)
	if entry.Phonetic != "" {
		fmt.Printf(" (%s)", entry.Phonetic)
	}
	fmt.Println()


	for _, meaning := range entry.Meanings {
		fmt.Printf("\n⚡ %s:\n", meaning.PartOfSpeech)

		for i, def := range meaning.Definitions {
			fmt.Printf("  %d. %s\n", i+1, def.Definition)

			if def.Example != "" {
				fmt.Printf("     🔸 Example: \"%s\"\n", def.Example)
			}

			if len(def.Synonyms) > 0 {
				fmt.Print("     🔁 Synonyms: ")
				for _, syn := range def.Synonyms {
					fmt.Printf("%s ", syn)
				}
				fmt.Println()
			}
		}
	}
}

