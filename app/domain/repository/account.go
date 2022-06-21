package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)

	// Fetch account which has specified ID
	FindByID(ctx context.Context, id int64) (*object.Account, error)

	// Create an account which has specified username and password
	CreateAccount(ctx context.Context, username, password string) error

	// Follow an account
	Follow(ctx context.Context, followerID, followeeID int64) (int64, bool, error)

	// Unfollow an account
	Unfollow(ctx context.Context, followerID, followeeID int64) (int64, bool, error)

	// Account relationship about follow
	FindRelationship(ctx context.Context, userID, targetID int64) (bool, bool, error)

	// Fetch accounts that followed by follower
	FindFollowing(ctx context.Context, followerID, limit int64) ([]object.Account, error)

	// Fetch accounts that following followee
	FindFollowers(ctx context.Context, followeeID, maxID, sinceID, limit int64) ([]object.Account, error)

	// Update credentials
	UpdateCredentials(ctx context.Context, id int64, displayName, note, avatar, header string) error
}
