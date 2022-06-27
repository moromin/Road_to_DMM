package repository

import (
	"context"
	"io"
	"yatter-backend-go/app/domain/object"
)

type Attachment interface {
	// Upload file
	UploadFile(ctx context.Context, file io.Reader, fileDir, filename, filetype, description string) (*object.Attachment, error)
}
