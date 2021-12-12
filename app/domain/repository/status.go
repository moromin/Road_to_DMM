package repository

import "context"

type Status interface {
	// Create a status which has specified content and authenticated account id
	CreateStatus(ctx context.Context, content string, accountID int64) error
}
