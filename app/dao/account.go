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
func (r *account) FindByID(ctx context.Context, id int) (*object.Account, error) {
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
func (r *account) Follow(ctx context.Context, follower_id, followee_id int64) (int64, bool, error) {
	var id int64
	var followedBy bool

	err := r.Transaction(func(tx *sqlx.Tx) error {
		const follow = `INSERT INTO follow (follower_id, followee_id) VALUES (?, ?)`
		res, err := tx.ExecContext(ctx, follow, follower_id, followee_id)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}

		if err := r.manageNumberOfFollows(ctx, tx, follower_id, "following_count", 1); err != nil {
			return err
		}
		if err := r.manageNumberOfFollows(ctx, tx, followee_id, "followers_count", 1); err != nil {
			return err
		}

		const followed = `SELECT * FROM follow WHERE follower_id = ? AND followee_id = ?`
		empty := struct{ I, J, K int64 }{}
		if err := r.db.QueryRowxContext(ctx, followed, followee_id, follower_id).Scan(&empty.I, &empty.J, &empty.K); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				followedBy = false
			} else {
				return err
			}
		} else {
			followedBy = true
		}

		return nil
	})
	if err != nil {
		return 0, false, err
	}

	return id, followedBy, nil
}

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

func (r *account) FindFollowing(ctx context.Context, follower_id, limit int64) ([]object.Account, error) {
	accounts := make([]object.Account, 0)

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

	for rows.Next() {
		account := object.Account{}
		err := rows.StructScan(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return accounts, nil
}
