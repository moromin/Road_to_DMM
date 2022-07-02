package mock

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

// AccountMock is a mock implementation of Account
type AccountMock struct {
	FindByUsernameFunc    func(ctx context.Context, username string) (*object.Account, error)
	FindByIDFunc          func(ctx context.Context, id int64) (*object.Account, error)
	CreateAccountFunc     func(ctx context.Context, username, password string) (int64, error)
	FollowFunc            func(ctx context.Context, followerID, followeeID int64) (int64, bool, error)
	UnfollowFunc          func(ctx context.Context, followerID, followeeID int64) (int64, bool, error)
	FindRelationshipFunc  func(ctx context.Context, userID, targetID int64) (bool, bool, error)
	FindFollowingFunc     func(ctx context.Context, followerID, limit int64) ([]object.Account, error)
	FindFollowersFunc     func(ctx context.Context, followeeID, maxID, sinceID, limit int64) ([]object.Account, error)
	UpdateCredentialsFunc func(ctx context.Context, id int64, displayName, note, avatar, header string) error
}

// FindByUsername is a mock implementation of Account.FindByUsername
func (m *AccountMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	return m.FindByUsernameFunc(ctx, username)
}

// FindByID is a mock implementation of Account.FindByID
func (m *AccountMock) FindByID(ctx context.Context, id int64) (*object.Account, error) {
	return m.FindByIDFunc(ctx, id)
}

// CreateAccount is a mock implementation of Account.CreateAccount
func (m *AccountMock) CreateAccount(ctx context.Context, username, password string) (int64, error) {
	return m.CreateAccountFunc(ctx, username, password)
}

// Follow is a mock implementation of Account.Follow
func (m *AccountMock) Follow(ctx context.Context, followerID, followeeID int64) (int64, bool, error) {
	return m.FollowFunc(ctx, followerID, followeeID)
}

// Unfollow is a mock implementation of Account.Unfollow
func (m *AccountMock) Unfollow(ctx context.Context, followerID, followeeID int64) (int64, bool, error) {
	return m.UnfollowFunc(ctx, followerID, followeeID)
}

// FindRelationship is a mock implementation of Account.FindRelationship
func (m *AccountMock) FindRelationship(ctx context.Context, userID, targetID int64) (bool, bool, error) {
	return m.FindRelationshipFunc(ctx, userID, targetID)
}

// FindFollowing is a mock implementation of Account.FindFollowing
func (m *AccountMock) FindFollowing(ctx context.Context, followerID, limit int64) ([]object.Account, error) {
	return m.FindFollowingFunc(ctx, followerID, limit)
}

// FindFollowers is a mock implementation of Account.FindFollowers
func (m *AccountMock) FindFollowers(ctx context.Context, followeeID, maxID, sinceID, limit int64) ([]object.Account, error) {
	return m.FindFollowersFunc(ctx, followeeID, maxID, sinceID, limit)
}

// UpdateCredentials is a mock implementation of Account.UpdateCredentials
func (m *AccountMock) UpdateCredentials(ctx context.Context, id int64, displayName, note, avatar, header string) error {
	return m.UpdateCredentialsFunc(ctx, id, displayName, note, avatar, header)
}
