package dao_test

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

// TODO: create test Main(), setup() and tearDown() to simplify tests
type DatabaseTestSuite struct {
	suite.Suite

	mockDB *sql.DB
	sqlxDB *sqlx.DB
	mock   sqlmock.Sqlmock
}

func (s *DatabaseTestSuite) setupDB() (*sql.DB, sqlmock.Sqlmock, *sqlx.DB) {
	mockDB, mock, err := sqlmock.New()
	s.Require().NoError(err)

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return mockDB, mock, sqlxDB
}

func (s *DatabaseTestSuite) setupSuite() {
	mockDB, mock, sqlxDB := s.setupDB()
	s.mockDB = mockDB
	s.mock = mock
	s.sqlxDB = sqlxDB
}

func (s *DatabaseTestSuite) tearDownSuite() {
	defer s.mockDB.Close()
	defer s.sqlxDB.Close()
}
