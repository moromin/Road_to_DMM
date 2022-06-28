package mock

import (
	"context"
	"io"
	"yatter-backend-go/app/domain/object"
)

// AttachmentMock is a mock implementation of Attachment
type AttachmentMock struct {
	UploadFileFunc func(ctx context.Context, file io.Reader, fileDir, filename, filetype, description string) (*object.Attachment, error)
}

// UploadFile is a mock implementation of Attachment.UploadFile
func (m *AttachmentMock) UploadFile(ctx context.Context, file io.Reader, fileDir, filename, filetype, description string) (*object.Attachment, error) {
	return m.UploadFileFunc(ctx, file, fileDir, filename, filetype, description)
}
