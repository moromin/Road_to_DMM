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
	"yatter-backend-go/app/domain/object"

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

// Get sqlx client
func dbClient() *sqlx.DB {
	sqlxDB := sqlx.NewDb(dockerDB.Conn, "mysql")
	dockerDB.Sqlx = sqlxDB
	return sqlxDB
}

// Make seeds for test
var (
	testUser1 = object.Account{ID: 1, Username: "test1", PasswordHash: "secret"}
	testUser2 = object.Account{ID: 2, Username: "test2", PasswordHash: "himitsu"}
	testUser3 = object.Account{ID: 3, Username: "test3", PasswordHash: "password"}
)

func makeSeeds(client *sqlx.DB) {
	// register accounts
	for _, user := range []object.Account{testUser1, testUser2, testUser3} {
		_, err := client.Exec(`INSERT INTO account (username, password_hash) VALUES (?, ?)`, user.Username, user.PasswordHash)
		if err != nil {
			log.Fatalf("Failed to register test accounts: %s", err)
		}
	}

	// follow accounts
	_, err := client.Exec(`INSERT INTO follow (follower_id, followee_id) VALUES (?, ?)`, 1, 2)
	if err != nil {
		log.Fatalf("Failed to follow test accounts: %s", err)
	}
	updateFollows := `UPDATE account SET %s = %s + %d WHERE id = %d`
	_, err = client.Exec(fmt.Sprintf(updateFollows, "following_count", "following_count", 1, 1))
	if err != nil {
		log.Fatalf("Failed to follow test accounts: %s", err)
	}
	_, err = client.Exec(fmt.Sprintf(updateFollows, "followers_count", "followers_count", 1, 2))
	if err != nil {
		log.Fatalf("Failed to follow test accounts: %s", err)
	}
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

	makeSeeds(dbClient())

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
		Message: fmt.Sprintf("Duplicate entry '%s' for key 'username'", testUser1.Username),
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
				id:  4,
				err: nil,
			},
		},
		"duplicate": {
			args: args{
				name:     testUser1.Username,
				password: testUser1.PasswordHash,
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

func TestAccount_FindByID(t *testing.T) {
	type args struct {
		id int64
	}
	type want struct {
		account *object.Account
		err     error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"success": {
			args: args{
				id: testUser1.ID,
			},
			want: want{
				account: &object.Account{
					ID:             testUser1.ID,
					Username:       testUser1.Username,
					PasswordHash:   testUser1.PasswordHash,
					FollowersCount: 0,
					FollowingCount: 1,
				},
				err: nil,
			},
		},
		"not found": {
			args: args{
				id: 0,
			},
			want: want{
				account: nil,
				err:     nil,
			},
		},
	}

	client := dbClient()
	repo := dao.NewAccount(client)
	ctx := context.Background()

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := repo.FindByID(ctx, tt.args.id)

			if tt.want.account != nil {
				assert.Equal(t, tt.want.account.ID, got.ID)
				assert.Equal(t, tt.want.account.PasswordHash, got.PasswordHash)
				assert.Equal(t, tt.want.account.Username, got.Username)
				assert.Equal(t, tt.want.account.FollowersCount, got.FollowersCount)
				assert.Equal(t, tt.want.account.FollowingCount, got.FollowingCount)
			} else {
				assert.Nil(t, got)
			}
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestAccount_Follow(t *testing.T) {
	type args struct {
		followerID int64
		followeeID int64
	}
	type want struct {
		id         int64
		followedBy bool
		err        error
	}

	ErrAlreadyFollowed := &mysql.MySQLError{
		Number:  0x426,
		Message: "Duplicate entry '1-2' for key 'follow_combination'",
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"success": {
			args: args{
				followerID: testUser3.ID,
				followeeID: testUser2.ID,
			},
			want: want{
				id:         2,
				followedBy: false,
				err:        nil,
			},
		},
		"already following": {
			args: args{
				followerID: testUser1.ID,
				followeeID: testUser2.ID,
			},
			want: want{
				id:         0,
				followedBy: false,
				err:        ErrAlreadyFollowed,
			},
		},
	}
	client := dbClient()
	repo := dao.NewAccount(client)
	ctx := context.Background()

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			id, followedBy, err := repo.Follow(ctx, tt.args.followerID, tt.args.followeeID)
			assert.Equal(t, tt.want.id, id)
			assert.Equal(t, tt.want.followedBy, followedBy)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
