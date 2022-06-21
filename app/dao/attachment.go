package dao

import (
	"context"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Attachment
	attachment struct {
		db *sqlx.DB
	}
)

func NewAttachment(db *sqlx.DB) repository.Attachment {
	return &attachment{db: db}
}

func (r *attachment) UploadFile(ctx context.Context) (int64, string, string, error)
