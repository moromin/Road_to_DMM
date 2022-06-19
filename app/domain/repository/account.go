package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)

	// Fetch account which has specified ID
	FindByID(ctx context.Context, id int) (*object.Account, error)

	// Create an account which has specified username and password
	CreateAccount(ctx context.Context, username, password string) error

	// Follow an account
	Follow(ctx context.Context, follower_id, followee_id int64) (int64, bool, error)

	// Fetch accounts that followed by follower
	FindFollowing(ctx context.Context, follower_id, limit int64) ([]object.Account, error)
}
