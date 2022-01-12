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

// Implementation for repository.Status
type status struct {
	db *sqlx.DB
}

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// CreateAccount : content, accountIDから新しいステータスを作成
func (r *status) Create(ctx context.Context, content string, accountID int64) error {
	query := `insert into status (
				account_id,
				content
			  ) values (
				?, ?
			  )`

	_, err := r.db.QueryContext(ctx, query, accountID, content)
	if err != nil {
		return err
	}
	return nil
}

// FindByID : IDからステータスを取得
func (r *status) FindByID(ctx context.Context, id int64) (*object.Status, error) {
	entity := new(object.Status)

	query := `select * 
			  from status
			  where id = ?`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}

// DeleteByID : IDからステータスを削除
func (r *status) DeleteByID(ctx context.Context, id int64) error {
	query := `delete 
			  from status
			  where id = ?`

	_, err := r.db.QueryxContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%w", err)
	}

	return nil
}

// List : maxID, sinceID, limit からタイムライン（ステータスのスライス）を取得
func (r *status) List(ctx context.Context, maxID, sinceID, limit int64) ([]object.Status, error) {
	where := ""
	max := ""
	since := ""
	and := ""

	fmt.Println(limit)
	if maxID != 0 || sinceID != 0 {
		where = "WHERE"
		if maxID != 0 {
			max = fmt.Sprintf("id <= %d", maxID)
		}
		if sinceID != 0 {
			since = fmt.Sprintf("id >= %d", sinceID)
		}
	}
	if maxID != 0 && sinceID != 0 {
		and = "AND"
	}

	query := fmt.Sprintf(`SELECT * FROM status %s %s %s %s LIMIT %d`, where, max, and, since, limit)

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
			log.Fatal(err)
		}
		statusList = append(statusList, status)
	}

	return statusList, nil
}
