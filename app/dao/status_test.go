package dao_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"yatter-backend-go/app/dao"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	SELECT = iota + 1
	INSERT
	DELETE
	UPDATE
)

// TODO: create test Main(), setup() and tearDown() to simplify tests

func Test_status_Create(t *testing.T) {
	// TDD

	const wantErr, noErr = true, false
	type args struct {
		ID      int64
		Content string
	}

	// TODO: export query strings
	cases := map[string]struct {
		query     string
		args      args
		expectErr bool
	}{
		"normal": {"INSERT INTO status ( account_id, content ) VALUES ( ?, ? )", args{1, "hogehoge"}, noErr},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockDB, mock, sqlxDB := MockDB(t)
			defer mockDB.Close()

			// TODO: mockClosure
			mockClosure(t, mock, INSERT, tt.query, tt.expectErr, convertStructToSlice(tt.args)...)
			// mock.ExpectExec(regexp.QuoteMeta(tt.query)).
			// 	WithArgs(tt.args.ID, tt.args.Content).
			// 	WillReturnResult(sqlmock.NewResult(1, 1))

			// TODO: how to handle context in test
			ctx := context.TODO()
			repo := dao.NewStatus(sqlxDB)

			if _, err := repo.Create(ctx, tt.args.ID, tt.args.Content, []int64{}); !isExpectedError(tt.expectErr, err) {
				t.Errorf("error was not expected: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func mockClosure(
	t *testing.T,
	mock sqlmock.Sqlmock,
	queryType int,
	query string,
	expectErr bool,
	args ...driver.Value,
) {
	t.Helper()
	switch queryType {
	case INSERT:
		ex := mock.ExpectExec(regexp.QuoteMeta(query))
		if expectErr {
			// TODO: handle error case
		} else {
			ex.WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
	}
}
