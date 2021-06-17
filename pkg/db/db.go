package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DriverName     = "sqlite3"
	DataSourceName = "perfectward.db"
)

type Database struct {
	Conn *sql.DB
}

// InitDB initialises and creates a connection to our database
func InitDB() (*Database, error) {
	// Initialise the database
	db, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
		return nil, err
	}
	log.Println("Successfully connected to the database!")

	// Test DB is ready to accept connections
	if err := db.Ping(); err != nil {
		log.Fatalf("DB is unable to accept connections: %v", err)
		return nil, err
	}
	log.Println("Ready to accept connections!")

	return &Database{Conn: db}, nil
}

// GetAllRows retrieves all rows from the averages table
func (db *Database) GetAllRows() (*sql.Rows, error) {
	numberOfRows, err := db.Conn.Query("SELECT * FROM averages;")
	if err != nil {
		return nil, err
	}

	return numberOfRows, nil
}

// CreateANewEntry creates a new entry based on the current run's values
func (db *Database) CreateANewEntry(numOfQuestions, positives int) error {
	stmt, err := db.Conn.Prepare("INSERT INTO averages(questions, positives) values(?,?);")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(numOfQuestions, positives)
	if err != nil {
		return err
	}

	return nil
}

// GetPersistedParams retrieves the persisted entries from db
func (db *Database) GetPersistedParams() (int, int, error) {
	var questions, positives int

	rows, err := db.Conn.Query("SELECT questions, positives FROM averages;")
	if err != nil {
		log.Fatal(err)
		return 0, 0, err
	}

	for rows.Next() {
		if err = rows.Scan(&questions, &positives); err != nil {
			log.Fatal(err)
			return 0, 0, err
		}
	}
	rows.Close()

	return questions, positives, nil
}

// UpdateDatabaseParams updates the table with new average
func (db *Database) UpdateDatabaseParams(totalQuestions, totalPositives int) error {
	// Replace the old average with new average
	stmt, err := db.Conn.Prepare("UPDATE averages SET questions=?, positives=? where uuid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(totalQuestions, totalPositives, 1)
	if err != nil {
		return err
	}

	return nil
}

// IsEmpty checks if there are any entries in the table
func (db *Database) IsEmpty(rows *sql.Rows) bool {
	if !rows.Next() {
		return true
	}
	rows.Close()

	return false
}
