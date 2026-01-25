// Package db handles the database connection and initialization.
package db

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Bahaaio/pomo/config"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const DBFile = config.AppName + ".db"

// Connect connects to the SQLite database,
// creates the necessary directories,
// and performs migrations if needed.
func Connect() (*sqlx.DB, error) {
	dbDir, err := getDBDir()
	if err != nil {
		log.Println("failed to get db path:", err)
		return nil, err
	}

	// create the db directory if it doesn't exist
	if err = os.MkdirAll(dbDir, 0o755); err != nil {
		log.Println("failed to create db directory:", err)
		return nil, err
	}

	dbPath := filepath.Join(dbDir, DBFile)

	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		log.Println("failed to connect to the db:", err)
		return nil, err
	}
	log.Println("connected to the db")

	if err = db.Ping(); err != nil {
		log.Println("failed to ping the db:", err)
		return nil, err
	}
	log.Println("pinged the db")

	// limit the number of open connections to 1
	db.SetMaxOpenConns(1)

	// migrate the database
	if err = createSchema(db); err != nil {
		log.Println("failed to migrate the db:", err)
		return nil, err
	}

	return db, nil
}

func createSchema(db *sqlx.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	log.Println("created the schema")

	return nil
}

// returns the path to the db directory
func getDBDir() (string, error) {
	var dir string

	// on Linux and macOS, use ~/.local/state
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}

		dir = filepath.Join(dir, ".local", "state")
	} else {
		// on other OSes, use the standard user config directory
		var err error
		dir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}

	// join the dir with the app name
	return filepath.Join(dir, config.AppName), nil
}
