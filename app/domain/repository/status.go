package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Create a status which has specified content and authenticated account id
	Create(ctx context.Context, content string, accountID int64) error

	// Fetch status which has specified ID
	FindByID(ctx context.Context, id int) (*object.Status, error)

	// Delete status which has specified ID
	DeleteByID(ctx context.Context, id int) error

	// Fetch statuses which has specified ID
	List(ctx context.Context, maxID, sinceID, limit int) ([]object.Status, error)
}
