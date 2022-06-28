package dao_test

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func isExpectedError(expect bool, actual error) bool {
	return expect == (actual != nil)
}

func convertStructToSlice(t *testing.T, x interface{}) []driver.Value {
	t.Helper()
	v := reflect.ValueOf(x)
	values := make([]driver.Value, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	return values
}

func convertOneToManyTable(t *testing.T, one interface{}, many ...interface{}) []driver.Value {
	t.Helper()
	values := make([]driver.Value, 2*len(many))
	for i, v := range many {
		values[2*i] = one
		values[2*i+1] = v
	}
	return values
}
