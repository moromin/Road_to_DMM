package dao

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
	"yatter-backend-go/app/domain/object"
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

const FilesPath = "localhost:8080/v1/media/"

func (r *attachment) UploadFile(ctx context.Context, file io.Reader, fileDir, filename, filetype, description string) (*object.Attachment, error) {
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return nil, err
	}

	dstName := fmt.Sprintf("%s/%d_%s", fileDir, time.Now().UnixNano(), filename)
	dst, err := os.Create(dstName)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	url := FilesPath + dstName

	const query = `INSERT INTO attachment (type, url, description) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, filetype, url, description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &object.Attachment{
		ID:          id,
		Type:        filetype,
		URL:         url,
		Description: description,
	}, nil
}
