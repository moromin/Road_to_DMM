package dao_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
)

// dockertest helpers
func createContainer(t *testing.T) (*dockertest.Resource, *dockertest.Pool) {
	t.Helper()

	pwd, _ := os.Getwd()

	pwd = pwd[:strings.Index(pwd, "/app")]

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	pool.MaxWait = time.Minute * 2

	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
		},
		Mounts: []string{
			pwd + "/.data/mysql:/etc/mysql",
			pwd + "/ddl:/docker-entrypoint-initdb.d",
		},
		Cmd: []string{
			"mysqld",
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_bin",
			"--default-time-zone='+9:00'",
		},
	}

	resource, err := pool.RunWithOptions(runOptions)
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	return resource, pool
}

func closeContainer(t *testing.T, resource *dockertest.Resource, pool *dockertest.Pool) {
	t.Helper()

	if err := pool.Purge(resource); err != nil {
		t.Fatalf("Could not purge resource: %s", err)
	}
}

func connectDB(t *testing.T, resource *dockertest.Resource, pool *dockertest.Pool) *sql.DB {
	t.Helper()

	var db *sql.DB
	// cfg := dbConfig()

	dsn := fmt.Sprintf("root:secret@(localhost:%s)/test_db?parseTime=true", resource.GetPort("3306/tcp"))

	if err := pool.Retry(func() error {
		time.Sleep(time.Second * 10)

		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return db
}

// utils
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
