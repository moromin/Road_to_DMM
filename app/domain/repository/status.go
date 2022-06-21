package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Create a status which has specified content and authenticated account id
	Create(ctx context.Context, id int64, content string) (int64, error)

	// Fetch status which has specified ID
	FindByID(ctx context.Context, id int64) (*object.Status, error)

	// Delete status which has specified ID
	DeleteByID(ctx context.Context, id int64) error

	// Fetch statuses
	ListAll(ctx context.Context, maxID, sinceID, limit int64) ([]object.Status, error)

	// Fetch statuses which has specified ID
	ListByID(ctx context.Context, id, maxID, sinceID, limit int64) ([]object.Status, error)
}
