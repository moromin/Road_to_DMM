package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/mock"
	"yatter-backend-go/app/domain/object"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestAccount_Create(t *testing.T) {
	type invalidRequest struct {
		Invalid string `json:"invalid"`
	}

	type args struct {
		*CreateRequest
		invalidRequest
	}
	type want struct {
		account *object.Account
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

			h := &handler{
				app: &app.App{Dao: dao.NewMock(
					&mock.AccountMock{
						CreateAccountFunc: func(ctx context.Context, username string, passwordHash string) error {
							return tt.want.err
						},
						FindByUsernameFunc: func(ctx context.Context, username string) (*object.Account, error) {
							return tt.want.account, tt.want.err
						},
					},
					nil,
					nil,
				)},
				validator: validator.New(),
			}

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

// t := s.T()

// 	resp, err := s.PostJSON("/v1/accounts", `{"name": "test", "password": "secret"}`)

// 	s.Assert().NoError(err)
// 	if !s.Assert().Equal(resp.StatusCode, http.StatusOK) {
// 		t.Skip("Status is not matched")
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if !s.Assert().NoError(err) {
// 		t.Skip("Failed to read body")
// 	}

// 	var j map[string]interface{}
// 	if s.Assert().NoError(json.Unmarshal(body, &j)) {
// 		s.Assert().Equal("test", j["name"])
// 	}
