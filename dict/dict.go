package dict

import (
    "fmt"
		"regexp"
		"strings"
)

type Definition struct {
    Definition string   `json:"definition"`
    Example    string   `json:"example"`
    Synonyms   []string `json:"synonyms"`
}

type Meaning struct {
    PartOfSpeech string       `json:"partOfSpeech"`
    Definitions  []Definition `json:"definitions"`
}

type DictionaryEntry struct {
    Word     string    `json:"word"`
    Phonetic string    `json:"phonetic"`
    Meanings []Meaning `json:"meanings"`
}

// Offline lookup clean
var (
    tagRe = regexp.MustCompile(`<[^>]+>`)
    splitRe = regexp.MustCompile(`(?i)<p>|<sn>\d+\.?</sn>`)
)

func cleanDefinition(raw string) string {
    // Split by <p> or numbered sense tags
    parts := splitRe.Split(raw, -1)
    var cleanedParts []string
    for _, part := range parts {
        part = strings.TrimSpace(part)
        if part == "" {
            continue
        }
        // Remove any remaining tags
        clean := tagRe.ReplaceAllString(part, "")
        cleanedParts = append(cleanedParts, clean)
    }
    return strings.Join(cleanedParts, "\n\n") // double newline between defs
}

func PrintEntry(entry *DictionaryEntry) {
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
