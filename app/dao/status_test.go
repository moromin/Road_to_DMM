package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type StatusTestSuite struct {
	DatabaseTestSuite

	repo repository.Status
}

func (s *StatusTestSuite) SetupSuite() {
	s.T().Log("SetupSuite")
	s.setupSuite()

	s.repo = dao.NewStatus(s.sqlxDB)
}

func (s *StatusTestSuite) TearDownSuite() {
	s.T().Log("TearDownSuite")
	s.tearDownSuite()
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusTestSuite))
}

func (s *StatusTestSuite) TestCreate() {
	type in struct {
		ID            int64
		Content       string
		AttachmentIDs []int64
	}
	type out struct {
		ID int64
	}
	const wantErr, noErr = true, false

	cases := map[string]struct {
		in        in
		want      out
		expectErr bool
	}{
		"Success with no attachment": {in{1, "simple", nil}, out{1}, noErr},
		// "Success with attachment": {in{1, "with attachments", []int64{1, 2}}, out{2}, noErr},
		// "Invalid attachment": {in{1, "invalid attachment", []int64{-1, 1}}, out{0}, wantErr},
	}

	t := s.T()
	ctx := context.Background()
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s.mock.ExpectBegin()
			s.mock.ExpectExec(`INSERT INTO status`).
				WithArgs(tt.in.ID, tt.in.Content).
				WillReturnResult(sqlmock.NewResult(1, 1))
			// if tt.in.AttachmentIDs != nil {
			// 	rows := sqlmock.NewRows([]string{"count"}).AddRow(len(tt.in.AttachmentIDs))
			// 	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM attachment WHERE id IN (?` + strings.Repeat(",?", len(tt.in.AttachmentIDs)-1) + `)`)).
			// 		WithArgs(tt.in.AttachmentIDs).
			// 		WillReturnRows(rows)
			// 	s.mock.ExpectExec(`INSERT INTO status_attachment`).
			// 		WithArgs(convertOneToManyTable(t, tt.in.ID, tt.in.AttachmentIDs)).
			// 		WillReturnResult(sqlmock.NewResult(int64(len(tt.in.AttachmentIDs)), int64(len(tt.in.AttachmentIDs))))
			// }
			if tt.expectErr {
				s.mock.ExpectRollback()
			} else {
				s.mock.ExpectCommit()
			}

			id, err := s.repo.Create(ctx, tt.in.ID, tt.in.Content, tt.in.AttachmentIDs)
			if tt.expectErr {
				s.Assert().Errorf(err, "want error, but no error")
			} else {
				s.Assert().NoErrorf(err, "want no error, but error")
			}
			s.Assert().Equal(tt.want.ID, id)
			if err := s.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
