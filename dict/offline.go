
package dict

import (
    "database/sql"
    "fmt"
    "strings"
    "sync"

    _ "github.com/mattn/go-sqlite3"
)

var (
    db    *sql.DB
    once  sync.Once
    dbErr error
)

// InitDB initializes the SQLite DB connection once.
func InitDB(dbPath string) error {
    once.Do(func() {
        db, dbErr = sql.Open("sqlite3", dbPath)
        if dbErr != nil {
            return
        }
        dbErr = db.Ping()
    })
    return dbErr
}

// LookupOffline looks up the definition of a word from the SQLite DB.
func LookupOffline(word string) (*DictionaryEntry, error) {
    if db == nil {
        return nil, fmt.Errorf("database not initialized, call InitDB first")
    }
    var rawDef string
    err := db.QueryRow("SELECT definition FROM words WHERE word = ? LIMIT 1", strings.ToLower(word)).Scan(&rawDef)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    cleanedDef := cleanDefinition(rawDef)

    entry := &DictionaryEntry{
        Word: word,
        Meanings: []Meaning{
            {
                PartOfSpeech: "", // unknown for now
                Definitions: []Definition{
                    {
                        Definition: cleanedDef,
                    },
                },
            },
        },
    }
    return entry, nil
}

// CloseDB closes the database connection (call on program exit).
func CloseDB() error {
    if db != nil {
        return db.Close()
    }
    return nil
}

