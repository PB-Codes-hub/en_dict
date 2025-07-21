package main

import (
    "flag"
    "fmt"
    "log"

    "github.com/PB-Codes-hub/en_dict/dict"
)

func main() {
    err := dict.InitDB("./data/gcide_dict.sqlite3")
    if err != nil {
        log.Fatalf("Failed to open DB: %v", err)
    }
    defer dict.CloseDB()

    // Parse command line flags and arguments
    offline := flag.Bool("offline", false, "Use offline dictionary lookup")
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        fmt.Println("Usage: en_dict [--offline] <word>")
        return
    }
    word := args[0]

    var entry *dict.DictionaryEntry

    if *offline {
        entry, err = dict.LookupOffline(word)
    } else {
        entry, err = dict.LookupOnline(word)
    }
    if err != nil {
        log.Fatalf("Lookup error: %v", err)
    }
    if entry == nil {
        fmt.Println("No definition found.")
        return
    }
    dict.PrintEntry(entry)
}
