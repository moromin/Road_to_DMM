package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return entity, nil
}

// FindByUsername : IDからユーザを取得
func (r *account) FindByID(ctx context.Context, id int64) (*object.Account, error) {
	entity := new(object.Account)
	query := `select *
			  from account
			  where id = ?`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return entity, nil
}

// CreateAccount : username, passwordから新しいアカウントを作成
func (r *account) CreateAccount(ctx context.Context, username, password string) error {
	query := `insert into account (
				username,
				password_hash
			  ) values (
				  ?, ?
			  )`
	_, err := r.db.QueryContext(ctx, query, username, password)
	if err != nil {
		return err
	}

	return nil
}

// Follow : アカウントをフォロー
func (r *account) Follow(ctx context.Context, followerID, followeeID int64) (int64, bool, error) {
	var followedBy bool

	err := r.Transaction(func(tx *sqlx.Tx) error {
		var err error
		const follow = `INSERT INTO follow (follower_id, followee_id) VALUES (?, ?)`
		if _, err := tx.ExecContext(ctx, follow, followerID, followeeID); err != nil {
			return err
		}

		if err := r.manageNumberOfFollows(ctx, tx, followerID, "following_count", 1); err != nil {
			return err
		}
		if err := r.manageNumberOfFollows(ctx, tx, followeeID, "followers_count", 1); err != nil {
			return err
		}

		if followedBy, err = r.findRelationship(ctx, followeeID, followerID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, false, err
	}

	return followeeID, followedBy, nil
}

// Unfollow : アカウントのフォロー解除
func (r *account) Unfollow(ctx context.Context, followerID, followeeID int64) (int64, bool, error) {
	var followedBy bool
	var empty struct{ I, J, K int64 }

	err := r.Transaction(func(tx *sqlx.Tx) error {
		var err error
		const following = `SELECT * FROM follow WHERE follower_id = ? AND followee_id = ?`
		if err := r.db.QueryRowxContext(ctx, following, followerID, followeeID).Scan(&empty.I, &empty.J, &empty.K); err != nil {
			return err
		}

		const unfollow = `DELETE FROM follow WHERE follower_id = ? AND followee_id = ?`
		if _, err := tx.ExecContext(ctx, unfollow, followerID, followeeID); err != nil {
			return err
		}

		if err := r.manageNumberOfFollows(ctx, tx, followerID, "following_count", -1); err != nil {
			return err
		}
		if err := r.manageNumberOfFollows(ctx, tx, followeeID, "followers_count", -1); err != nil {
			return err
		}

		if followedBy, err = r.findRelationship(ctx, followeeID, followerID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, false, err
	}

	return followeeID, followedBy, nil
}

// Manage number of follower count and following count
func (r *account) manageNumberOfFollows(ctx context.Context, tx *sqlx.Tx, id int64, column string, number int64) error {
	var operator string
	if number == 0 {
		return nil
	} else if number > 0 {
		operator = "+"
	} else {
		operator = "-"
		number *= -1
	}

	query := fmt.Sprintf("UPDATE account SET %s = %s %s %d WHERE id = %d", column, column, operator, number, id)
	_, err := tx.ExecContext(ctx, query)

	return err
}

// Transaction handle specific process
// Essentially, it should be abstracted by DB interface
func (r *account) Transaction(txFunc func(*sqlx.Tx) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			log.Println("rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(tx)
	return err
}

// FindRelationship : 指定したアカウントとのフォロー関係を取得する
func (r *account) FindRelationship(ctx context.Context, userID, targetID int64) (bool, bool, error) {
	following, err := r.findRelationship(ctx, userID, targetID)
	if err != nil {
		return false, false, err
	}

	followedBy, err := r.findRelationship(ctx, targetID, userID)
	if err != nil {
		return false, false, err
	}

	return following, followedBy, nil
}

func (r *account) findRelationship(ctx context.Context, followerID, followeeID int64) (bool, error) {
	const query = `SELECT * FROM follow WHERE follower_id = ? AND followee_id = ?`
	var empty struct{ I, J, K int64 }

	if err := r.db.QueryRowxContext(ctx, query, followerID, followeeID).Scan(&empty.I, &empty.J, &empty.K); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// FindFollowing : フォローしているアカウント情報を取得する
func (r *account) FindFollowing(ctx context.Context, follower_id, limit int64) ([]object.Account, error) {
	query := `SELECT a.*
				FROM follow as f
				JOIN account as a
				ON f.followee_id = a.id
				WHERE follower_id = ?
				LIMIT ?`

	rows, err := r.db.QueryxContext(ctx, query, follower_id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]object.Account, 0)
	for rows.Next() {
		account := object.Account{}
		err := rows.StructScan(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindFollowers : フォローされているアカウント情報を取得する
func (r *account) FindFollowers(ctx context.Context, followeeID, maxID, sinceID, limit int64) ([]object.Account, error) {
	connection := ""
	idRange, ok := BuildRangeQuery("a.id", maxID, sinceID, 0)
	if ok {
		connection = "AND"
	} else {
		connection = "WHERE"
	}
	query := fmt.Sprintf(`SELECT a.*
								FROM follow as f
								JOIN account as a
								ON f.follower_id = a.id
								%s %s followee_id = %d
								LIMIT %d`, idRange, connection, followeeID, limit)

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]object.Account, 0)
	for rows.Next() {
		account := object.Account{}
		err := rows.StructScan(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// UpdateCredentials : アカウントの経歴を更新する
func (r *account) UpdateCredentials(ctx context.Context, id int64, displayName, note, avatar, header string) error {
	var columns string

	credentials := map[string]string{
		"display_name": displayName,
		"note":         note,
		"avatar":       avatar,
		"header":       header,
	}

	for name, value := range credentials {
		if value == "" {
			continue
		}
		if columns != "" {
			columns += ", "
		}
		columns += fmt.Sprintf("%s = %q", name, value)
	}
	if columns == "" {
		return nil
	}

	query := fmt.Sprintf("UPDATE account SET %s WHERE id = %d", columns, id)
	_, err := r.db.ExecContext(ctx, query)
	return err
}
