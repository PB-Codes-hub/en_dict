package dict

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)


func LookupOnline(word string) (*DictionaryEntry, error) {
    url := "https://api.dictionaryapi.dev/api/v2/entries/en/" + word
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("HTTP error: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("no definition found for: %s", word)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %v", err)
    }

    var results []DictionaryEntry
    if err := json.Unmarshal(body, &results); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }

    return &results[0], nil // return the first result
}
