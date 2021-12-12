package dao

import (
	"context"
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
func (r *status) CreateStatus(ctx context.Context, content string, accountID int64) error {
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
