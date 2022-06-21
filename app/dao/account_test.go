package dao

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"
	"yatter-backend-go/app/domain/object"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_account_FindByUsername(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}

	query := "select * from account where username = ?"
	displayName := "tester"
	avater := "none"
	header := "null"
	note := "nothing to write"

	want := &object.Account{
		ID:           1,
		Username:     "test user",
		PasswordHash: "sdfadgadsgaaga",
		DisplayName:  &displayName,
		Avatar:       &avater,
		Header:       &header,
		Note:         &note,
		CreateAt:     object.DateTime{Time: time.Now()},
	}

	tests := []struct {
		name        string
		mockClosure func(sqlmock.Sqlmock)
		args        args
		want        *object.Account
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "display_name", "avatar", "header", "note", "create_at"}).
					AddRow(want.ID, want.Username, want.PasswordHash, want.DisplayName,
						want.Avatar, want.Header, want.Note, want.CreateAt)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(want.Username).
					WillReturnRows(rows)
			},
			args: args{
				ctx:      context.Background(),
				username: want.Username,
			},
			want:      want,
			assertion: assert.NoError,
		},
		{
			name: "Failure to select",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(fmt.Errorf("select error"))
			},
			args: args{
				ctx:      context.Background(),
				username: want.Username,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("failed to init db mock:", err)
			}
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			defer sqlxDB.Close()

			tt.mockClosure(mock)

			r := &account{
				db: sqlxDB,
			}
			got, err := r.FindByUsername(tt.args.ctx, tt.args.username)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_account_FindByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}

	query := `select * from account where id = ?`
	displayName := "tester"
	avater := "none"
	header := "null"
	note := "nothing to write"

	want := &object.Account{
		ID:           1,
		Username:     "test user",
		PasswordHash: "sdfadgadsgaaga",
		DisplayName:  &displayName,
		Avatar:       &avater,
		Header:       &header,
		Note:         &note,
		CreateAt:     object.DateTime{Time: time.Now()},
	}

	tests := []struct {
		name        string
		mockClosure func(sqlmock.Sqlmock)
		args        args
		want        *object.Account
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "display_name", "avatar", "header", "note", "create_at"}).
					AddRow(want.ID, want.Username, want.PasswordHash, want.DisplayName,
						want.Avatar, want.Header, want.Note, want.CreateAt)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(want.ID).
					WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
				id:  want.ID,
			},
			want:      want,
			assertion: assert.NoError,
		},
		{
			name: "Failure to select",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(fmt.Errorf("select error"))
			},
			args: args{
				ctx: context.Background(),
				id:  want.ID,
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("failed to init db mock:", err)
			}
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			defer sqlxDB.Close()

			tt.mockClosure(mock)

			r := &account{
				db: sqlxDB,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_account_CreateAccount(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
		password string
	}

	query := `insert into account ( username, password_hash ) values ( ?, ? )`
	username := "tester"
	password_hash := "dfaewafasdfaefa"

	tests := []struct {
		name        string
		mockClosure func(sqlmock.Sqlmock)
		args        args
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password_hash"}).
					AddRow(1, username, password_hash)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(rows)
			},
			args: args{
				ctx:      context.Background(),
				username: username,
				password: password_hash,
			},
			assertion: assert.NoError,
		},
		{
			name: "Failure to insert",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(username, password_hash).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("error"))).
					WillReturnError(errors.New("insert failed"))
			},
			args: args{
				ctx:      context.Background(),
				username: username,
				password: password_hash,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				log.Fatal("failed to init db mock:", err)
			}
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			defer sqlxDB.Close()

			tt.mockClosure(mock)

			r := &account{
				db: sqlxDB,
			}
			tt.assertion(t, r.CreateAccount(tt.args.ctx, tt.args.username, tt.args.password))
		})
	}
}
