package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

// Implementation for repository.Status
type status struct {
	db *sqlx.DB
}

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// CreateAccount : content, accountIDから新しいステータスを作成
func (r *status) Create(ctx context.Context, accountID int64, content string) (int64, error) {
	query := `INSERT INTO status (
				account_id,
				content
			  ) VALUES (
				?, ?
			  )`

	// _, err := r.db.QueryContext(ctx, query, accountID, content)
	res, err := r.db.ExecContext(ctx, query, accountID, content)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// FindByID : IDからステータスを取得
func (r *status) FindByID(ctx context.Context, id int64) (*object.Status, error) {
	entity := &object.Status{}

	query := `select *
			  from status
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

// DeleteByID : IDからステータスを削除
func (r *status) DeleteByID(ctx context.Context, id int64) error {
	query := `delete
			  from status
			  where id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("status %d is not found", id)
	}

	return nil
}

// ListAll : maxID, sinceID, limit からタイムライン（ステータスのスライス）を取得
func (r *status) ListAll(ctx context.Context, maxID, sinceID, limit int64) ([]object.Status, error) {
	idRange, _ := BuildRangeQuery("s.id", maxID, sinceID, 0)
	query := fmt.Sprintf(`SELECT s.*, a.username AS "account.username", a.followers_count AS "account.followers_count", a.following_count AS "account.following_count", a.create_at AS "account.create_at"
							FROM status as s
							JOIN account as a
							on s.account_id = a.id
							%s
							ORDER BY s.id
							LIMIT %d`, idRange, limit)

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var statusList []object.Status
	var status object.Status
	for rows.Next() {
		err := rows.StructScan(&status)
		if err != nil {
			return nil, err
		}
		statusList = append(statusList, status)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return statusList, nil
}

// ListByID : 認証されたアカウントのID, maxID, sinceID, limit からタイムライン（ステータスのスライス）を取得
func (r *status) ListByID(ctx context.Context, id, maxID, sinceID, limit int64) ([]object.Status, error) {
	connection := ""
	idRange, ok := BuildRangeQuery("id", maxID, sinceID, 0)
	if ok {
		connection = "AND"
	} else {
		connection = "WHERE"
	}
	query := fmt.Sprintf(`SELECT * FROM status %s %s account_id = %d LIMIT %d`, idRange, connection, id, limit)
	fmt.Println(query)
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}

	var statusList []object.Status
	var status object.Status
	for rows.Next() {
		err := rows.StructScan(&status)
		if err != nil {
			return nil, err
		}
		statusList = append(statusList, status)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return statusList, nil
}
