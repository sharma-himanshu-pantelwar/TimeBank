package persistance

import (
	"database/sql"
	"djson/internal/config"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB // pointer to sql.DB
}

// Function to connect to database
// a function which accepts no parameters and returns a pointer to DB and error (if any)
func NewDatabase() (*Database, error) {
	// Load the config struct from config package
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME, config.DB_SSLMODE)

	fmt.Println("Database URL: ", databaseURL)

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

// Function to close database
func (d *Database) Close() {
	d.db.Close()
}

// Func to get Db
func (d *Database) GetDB() *sql.DB {
	return d.db
}
