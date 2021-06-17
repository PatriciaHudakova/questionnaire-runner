//go:generate mockgen -source=interface.go -destination=mock.go -package db
package db

import "database/sql"

type IDatabase interface {
	GetAllRows() (*sql.Rows, error)
	CreateANewEntry(numOfQuestions, positives int) error
	GetPersistedParams() (int, int, error)
	UpdateAverage(totalQuestions, totalPositives int) error
	IsEmpty(rows *sql.Rows) bool
}
