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
func (r *status) FindByID(ctx context.Context, id int) (*object.Status, error) {
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
