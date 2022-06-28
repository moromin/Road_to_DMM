package testutils

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
	"time"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

// Test suite
type Suite struct {
	suite.Suite

	Resource *dockertest.Resource
	Pool     *dockertest.Pool
	Conn     *sql.DB
	SqlxDB   *sqlx.DB
	App      *app.App
	Server   *httptest.Server
}

// Create docker environment for test
func CreateContainer(t *testing.T) (*dockertest.Resource, *dockertest.Pool) {
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

func CloseContainer(t *testing.T, resource *dockertest.Resource, pool *dockertest.Pool) {
	t.Helper()

	if err := pool.Purge(resource); err != nil {
		t.Fatalf("Could not purge resource: %s", err)
	}
}

// func dbConfig() *mysql.Config {
// 	cfg := mysql.NewConfig()

// 	cfg.User = "root"
// 	cfg.Passwd = "secret"
// 	cfg.Net = "tcp"
// 	cfg.Addr = "localhost:3306"
// 	cfg.DBName = "mydb"
// 	cfg.ParseTime = true
// 	cfg.Loc = time.FixedZone("Asia/Tokyo", 9*60*60)

// 	return cfg
// }

func ConnectDB(t *testing.T, resource *dockertest.Resource, pool *dockertest.Pool) *sql.DB {
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

// Common test helper
func (s *Suite) Setup() {
	t := s.T()
	s.Resource, s.Pool = CreateContainer(t)
	s.Conn = ConnectDB(t, s.Resource, s.Pool)
	s.SqlxDB = sqlx.NewDb(s.Conn, "mysql")
	s.App = &app.App{Dao: dao.NewDao(s.SqlxDB)}
}

func (s *Suite) Teardown() {
	t := s.T()
	s.Server.Close()
	CloseContainer(t, s.Resource, s.Pool)
}

// Handler test helper
func (s *Suite) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return s.Server.Client().Post(s.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

func (s *Suite) Get(apiPath string) (*http.Response, error) {
	return s.Server.Client().Get(s.asURL(apiPath))
}

func (s *Suite) asURL(apiPath string) string {
	baseURL, _ := url.Parse(s.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
