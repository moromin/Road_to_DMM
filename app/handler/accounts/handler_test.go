package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/mock"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

const (
	wantErr = true
	noErr   = false
)

func TestAccount_Create(t *testing.T) {
	t.Parallel()

	type invalidRequest struct {
		Invalid string `json:"invalid"`
	}

	type args struct {
		*CreateRequest
		invalidRequest
	}
	type want struct {
		account *object.Account
		id      int64
		status  int
		err     error
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"success": {
			args: args{
				CreateRequest: &CreateRequest{
					Username: "test",
					Password: "secret",
				},
			},
			want: want{
				account: &object.Account{
					Username:       "test",
					FollowersCount: 0,
					FollowingCount: 0,
				},
				id:     1,
				status: http.StatusCreated,
				err:    nil,
			},
		},
		"invalid request": {
			args: args{
				invalidRequest: invalidRequest{
					Invalid: "invalid",
				},
			},
			want: want{
				account: nil,
				status:  http.StatusBadRequest,
				err:     errors.New("invalid request"),
			},
		},
		"failed to create account": {
			args: args{
				CreateRequest: &CreateRequest{
					Username: "test",
					Password: "secret",
				},
			},
			want: want{
				account: nil,
				status:  http.StatusInternalServerError,
				err:     errors.New("failed to create account"),
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {

			var buf bytes.Buffer
			var err error
			if tt.args.CreateRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.CreateRequest)
			} else {
				err = json.NewEncoder(&buf).Encode(tt.args.invalidRequest)
			}
			assert.Nil(t, err)

			r := httptest.NewRequest(http.MethodPost, "/v1/accounts", &buf)
			w := httptest.NewRecorder()

			router := chi.NewRouter()

			app := &app.App{Dao: dao.NewMock(
				&mock.AccountMock{
					CreateAccountFunc: func(ctx context.Context, username string, passwordHash string) (int64, error) {
						return tt.want.id, tt.want.err
					},
					FindByIDFunc: func(ctx context.Context, id int64) (*object.Account, error) {
						return tt.want.account, tt.want.err
					},
				},
				nil,
				nil,
			)}

			v := validator.New()

			h, _ := newHandlerAndRouter(router, app, v)

			h.Create(w, r)

			assert.Equal(t, tt.want.status, w.Code)

			var got object.Account
			err = json.NewDecoder(w.Body).Decode(&got)
			if tt.want.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			if tt.want.account != nil {
				assert.Equal(t, tt.want.account.Username, got.Username)
				assert.Equal(t, tt.want.account.FollowersCount, got.FollowersCount)
				assert.Equal(t, tt.want.account.FollowingCount, got.FollowingCount)
			}
		})
	}
}

func TestAccount_Get(t *testing.T) {
	t.Parallel()

	var (
		displayName = "hoger"
		avatar      = "blue.png"
		header      = "green.png"
		note        = "This is a note"
	)

	type args struct {
		username string
	}
	type want struct {
		account *object.Account
		status  int
		err     error
	}
	cases := map[string]struct {
		args      args
		want      want
		expectErr bool
	}{
		"success": {
			args: args{
				username: "test",
			},
			want: want{
				account: &object.Account{
					Username:       "test",
					DisplayName:    &displayName,
					FollowersCount: 3,
					FollowingCount: 5,
					Avatar:         &avatar,
					Header:         &header,
					Note:           &note,
				},
				status: http.StatusOK,
				err:    nil,
			},
			expectErr: noErr,
		},
		"not found": {
			args: args{
				username: "no_one",
			},
			want: want{
				account: nil,
				status:  http.StatusNotFound,
				err:     nil,
			},
			expectErr: wantErr,
		},
		"failed to find account": {
			args: args{
				username: "test",
			},
			want: want{
				account: nil,
				status:  http.StatusInternalServerError,
				err:     errors.New("failed to find account"),
			},
			expectErr: wantErr,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/accounts/{%s}", tt.args.username), nil)
			w := httptest.NewRecorder()

			router := chi.NewRouter()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("username", tt.args.username)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			app := &app.App{Dao: dao.NewMock(
				&mock.AccountMock{
					FindByUsernameFunc: func(ctx context.Context, username string) (*object.Account, error) {
						return tt.want.account, tt.want.err
					},
				},
				nil,
				nil,
			)}

			v := validator.New()

			h, _ := newHandlerAndRouter(router, app, v)

			h.Get(w, r)

			assert.Equal(t, tt.want.status, w.Code)

			var err error
			var got object.Account
			err = json.NewDecoder(w.Body).Decode(&got)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			if tt.want.account != nil {
				assert.Equal(t, tt.want.account.Username, got.Username)
				assert.Equal(t, tt.want.account.FollowersCount, got.FollowersCount)
				assert.Equal(t, tt.want.account.FollowingCount, got.FollowingCount)
				assert.Equal(t, tt.want.account.DisplayName, got.DisplayName)
				assert.Equal(t, tt.want.account.Avatar, got.Avatar)
				assert.Equal(t, tt.want.account.Header, got.Header)
				assert.Equal(t, tt.want.account.Note, got.Note)
			}
		})
	}
}

func TestAccount_Follow(t *testing.T) {
	t.Parallel()

	follower := &object.Account{
		ID:       1,
		Username: "follower",
	}

	type args struct {
		username string
	}
	type want struct {
		account      *object.Account
		relationship *Relationship
		status       int
		err          error
	}
	cases := map[string]struct {
		args      args
		want      want
		expectErr bool
	}{
		"success": {
			args: args{
				username: "test",
			},
			want: want{
				account: &object.Account{
					ID:             2,
					Username:       "test",
					FollowersCount: 3,
					FollowingCount: 5,
				},
				relationship: &Relationship{
					ID:         1,
					Following:  true,
					FollowedBy: false,
				},
				status: http.StatusOK,
				err:    nil,
			},
			expectErr: noErr,
		},
		"not found": {
			args: args{
				username: "no_one",
			},
			want: want{
				account: nil,
				relationship: &Relationship{
					ID:         -1,
					Following:  false,
					FollowedBy: false,
				},
				status: http.StatusBadRequest,
				err:    nil,
			},
			expectErr: wantErr,
		},
		"follower yourself": {
			args: args{
				username: "follower",
			},
			want: want{
				account:      nil,
				relationship: nil,
				status:       http.StatusBadRequest,
				err:          nil,
			},
			expectErr: wantErr,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/accounts/{%s}/follow", tt.args.username), nil)
			w := httptest.NewRecorder()

			r = auth.SetAccount(r, follower)

			router := chi.NewRouter()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("username", tt.args.username)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			app := &app.App{Dao: dao.NewMock(
				&mock.AccountMock{
					FindByUsernameFunc: func(ctx context.Context, username string) (*object.Account, error) {
						return tt.want.account, tt.want.err
					},
					FollowFunc: func(ctx context.Context, followingID, followeeID int64) (int64, bool, error) {
						return tt.want.relationship.ID, tt.want.relationship.FollowedBy, tt.want.err
					},
				},
				nil,
				nil,
			)}
			v := validator.New()

			h, _ := newHandlerAndRouter(router, app, v)

			h.Follow(w, r)

			assert.Equal(t, tt.want.status, w.Code)

			var err error
			var got Relationship
			err = json.NewDecoder(w.Body).Decode(&got)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			if !tt.expectErr {
				assert.Equal(t, tt.want.relationship.ID, got.ID)
				assert.Equal(t, tt.want.relationship.Following, got.Following)
				assert.Equal(t, tt.want.relationship.FollowedBy, got.FollowedBy)
			}
		})
	}
}
