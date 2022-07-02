package dao_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
	"yatter-backend-go/app/dao"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	dockerDB *db
)

type db struct {
	Conn *sql.DB
	Sqlx *sqlx.DB
}

func dbClient() *sqlx.DB {
	sqlxDB := sqlx.NewDb(dockerDB.Conn, "mysql")
	dockerDB.Sqlx = sqlxDB
	return sqlxDB
}

func TestMain(m *testing.M) {
	// dockertest
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	pwd, _ := os.Getwd()

	pwd = pwd[:strings.Index(pwd, "/app")]

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_DATABASE=mydb",
		},
		Mounts: []string{
			pwd + "/ddl:/docker-entrypoint-initdb.d",
		},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("3306/tcp")
	databaseUrl := fmt.Sprintf("root:secret@tcp(%s)/mydb?parseTime=true", hostAndPort)
	log.Println(databaseUrl)

	_ = resource.Expire(120)

	dockerDB = &db{}

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		time.Sleep(time.Second * 10)

		dockerDB.Conn, err = sql.Open("mysql", databaseUrl)
		if err != nil {
			return err
		}
		dockerDB.Conn.SetConnMaxLifetime(time.Second)
		return dockerDB.Conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestAccount_Create(t *testing.T) {
	type args struct {
		name     string
		password string
	}
	type want struct {
		id  int64
		err error
	}

	const (
		testUsername = "test"
		testPassword = "secret"
	)

	ErrDupicateEntry := &mysql.MySQLError{
		Number:  1062,
		Message: fmt.Sprintf("Duplicate entry '%s' for key 'username'", testUsername),
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"success": {
			args: args{
				name:     testUsername,
				password: testPassword,
			},
			want: want{
				id:  1,
				err: nil,
			},
		},
		"duplicate": {
			args: args{
				name:     testUsername,
				password: testPassword,
			},
			want: want{
				id:  0,
				err: ErrDupicateEntry,
			},
		},
	}

	client := dbClient()
	repo := dao.NewAccount(client)
	ctx := context.Background()

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := repo.CreateAccount(ctx, tt.args.name, tt.args.password)
			assert.Equal(t, tt.want.id, got)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
