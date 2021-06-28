package db

import (
	"database/sql"
)

// MockClient mock.
type MockClient struct {
	GetAllRowsFn         func() (*sql.Rows, error)
	CreateANewEntryFn    func(numOfQuestions, positives int) error
	GetPersistedParamsFn func() (int, int, error)
	UpdateAverageFn      func(totalQuestions, totalPositives int) error
	IsEmptyFn            func(rows *sql.Rows) bool
}

// NewMockClient mock.
func NewMockClient() *MockClient {
	return &MockClient{
		GetAllRowsFn: func() (*sql.Rows, error) {
			return &sql.Rows{}, nil
		},
		CreateANewEntryFn: func(numOfQuestions, positives int) error {
			return nil
		},
		GetPersistedParamsFn: func() (int, int, error) {
			return 30, 2, nil
		},
		UpdateAverageFn: func(totalQuestions, totalPositives int) error {
			return nil
		},
		IsEmptyFn: func(rows *sql.Rows) bool {
			return true
		},
	}
}

// GetAllRows mock.
func (m *MockClient) GetAllRows() (*sql.Rows, error) {
	return m.GetAllRowsFn()
}

// CreateANewEntry mock.
func (m *MockClient) CreateANewEntry(numOfQuestions, positives int) error {
	return m.CreateANewEntryFn(numOfQuestions, positives)
}

// GetPersistedParams mock.
func (m *MockClient) GetPersistedParams() (int, int, error) {
	return m.GetPersistedParamsFn()
}

// UpdateAverage mock.
func (m *MockClient) UpdateDatabaseParams(totalQuestions, totalPositives int) error {
	return m.UpdateAverageFn(totalQuestions, totalPositives)
}

// IsEmpty mock.
func (m *MockClient) IsEmpty(rows *sql.Rows) bool {
	return m.IsEmptyFn(rows)
}
