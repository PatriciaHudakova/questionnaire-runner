package db

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(testConn *sql.DB) {
	// Drop averages table if existing in case of leftover resources from failed tests
	query, err := testConn.Prepare("DROP TABLE averages;")
	if err == nil {
		// Error is thrown if db exists, if not, drop it
		if _, err = query.Exec(); err != nil {
			log.Fatal(err)
		}
	}

	// Create the new averages table
	query, err = testConn.Prepare("CREATE TABLE averages (uuid INTEGER PRIMARY KEY, questions INTEGER, positives INTEGER);")
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func teardown(testConn *sql.DB) {
	testConn.Close()

	// Drop averages table at the end of successful test cases
	query, err := testConn.Prepare("DROP TABLE averages;")
	if err == nil {
		if _, err = query.Exec(); err != nil {
			log.Fatal(err)
		}
	}
}

func TestInitDB(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should be able to initialise db and connect to it",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := InitDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Assert the received db is not nil
			assert.NotNil(t, db)

			// Close the db connection during cleanup
			t.Cleanup(func() {
				db.Conn.Close()
			})
		})
	}
}

func TestDatabase_GetAllRows(t *testing.T) {
	testConn, _ := sql.Open(DriverName, DataSourceName)
	defer testConn.Close()

	DB := Database{Conn: testConn}

	// Setup
	setup(DB.Conn)
	// Cleanup
	defer teardown(DB.Conn)

	// Should not throw an error when the query is successful but no entries found
	rows, err := DB.GetAllRows()
	assert.Nil(t, err)
	assert.NotNil(t, rows)

	// Add an entry into the table
	entry, err := DB.Conn.Prepare("INSERT INTO averages (questions, positives) values(?,?)")
	if err != nil {
		t.Fatal()
	}
	res, err := entry.Exec(5, 3)
	if err != nil {
		t.Fatal()
	}

	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatal()
	}

	// Check one entry has been added
	assert.Equal(t, affected, int64(1))

	// Should not throw an error when searching for populated entries
	rows, err = DB.GetAllRows()
	assert.Nil(t, err)
	assert.NotNil(t, rows)
}

func TestDatabase_GetPersistedParams(t *testing.T) {
	testConn, _ := sql.Open(DriverName, DataSourceName)
	defer testConn.Close()

	DB := Database{Conn: testConn}

	// Setup
	setup(DB.Conn)
	// Cleanup
	defer teardown(DB.Conn)

	// Insert a mock average value into the table
	stmt, err := DB.Conn.Prepare("INSERT INTO averages(questions, positives) values(?,?);")
	if err != nil {
		t.Fatal()
	}

	_, err = stmt.Exec(10, 8)
	if err != nil {
		t.Fatal()
	}

	// Check that we are able to retrieve this average from table
	questions, positives, err := DB.GetPersistedParams()
	assert.Nil(t, err)
	assert.Equal(t, 10, questions)
	assert.Equal(t, 8, positives)
}

func TestDatabase_CreateANewEntry(t *testing.T) {
	testConn, _ := sql.Open(DriverName, DataSourceName)
	defer testConn.Close()

	DB := Database{Conn: testConn}

	// Setup
	setup(DB.Conn)
	// Cleanup
	defer teardown(DB.Conn)

	// Execution of query should not throw an error
	err := DB.CreateANewEntry(5, 3)
	assert.Nil(t, err)

	// Check the correct number has been populated
	questions, positives, err := DB.GetPersistedParams()
	assert.Nil(t, err)
	assert.Equal(t, 5, questions)
	assert.Equal(t, 3, positives)
}

func TestDatabase_UpdateDatabaseParams(t *testing.T) {
	testConn, _ := sql.Open(DriverName, DataSourceName)
	defer testConn.Close()

	DB := Database{Conn: testConn}

	// Setup
	setup(DB.Conn)
	// Cleanup
	defer teardown(DB.Conn)

	// Insert a mock existing entry
	err := DB.CreateANewEntry(5, 3)
	assert.Nil(t, err)

	// Update the mock entry with another entry
	err = DB.UpdateDatabaseParams(10, 4)
	assert.Nil(t, err)

	// Check the replacement has been successful
	questions, positives, err := DB.GetPersistedParams()
	assert.Nil(t, err)
	assert.Equal(t, 10, questions)
	assert.Equal(t, 4, positives)
}
