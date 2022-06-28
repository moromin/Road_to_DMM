package mock

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

// StatusMock is a mock implementation of Status
type StatusMock struct {
	CreateFunc     func(ctx context.Context, id int64, content string, attachmentIDs []int64) (int64, error)
	FindByIDFunc   func(ctx context.Context, id int64) (*object.Status, error)
	DeleteByIDFunc func(ctx context.Context, id int64) error
	ListAllFunc    func(ctx context.Context, maxID, sinceID, limit int64) ([]object.Status, error)
	ListByIDFunc   func(ctx context.Context, id, maxID, sinceID, limit int64) ([]object.Status, error)
}

// Create is a mock implementation of Status.Create
func (m *StatusMock) Create(ctx context.Context, id int64, content string, attachmentIDs []int64) (int64, error) {
	return m.CreateFunc(ctx, id, content, attachmentIDs)
}

// FindByID is a mock implementation of Status.FindByID
func (m *StatusMock) FindByID(ctx context.Context, id int64) (*object.Status, error) {
	return m.FindByIDFunc(ctx, id)
}

// DeleteByID is a mock implementation of Status.DeleteByID
func (m *StatusMock) DeleteByID(ctx context.Context, id int64) error {
	return m.DeleteByIDFunc(ctx, id)
}

// ListAll is a mock implementation of Status.ListAll
func (m *StatusMock) ListAll(ctx context.Context, maxID, sinceID, limit int64) ([]object.Status, error) {
	return m.ListAllFunc(ctx, maxID, sinceID, limit)
}

// ListByID is a mock implementation of Status.ListByID
func (m *StatusMock) ListByID(ctx context.Context, id, maxID, sinceID, limit int64) ([]object.Status, error) {
	return m.ListByIDFunc(ctx, id, maxID, sinceID, limit)
}
