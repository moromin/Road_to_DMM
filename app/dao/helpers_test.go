package dao_test

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func MockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlx.DB) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to initialize mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return mockDB, mock, sqlxDB
}

func isExpectedError(expect bool, actual error) bool {
	return expect == (actual != nil)
}

func convertStructToSlice(x interface{}) []driver.Value {
	v := reflect.ValueOf(x)
	values := make([]driver.Value, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	return values
}
