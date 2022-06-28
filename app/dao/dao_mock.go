package dao

import (
	"yatter-backend-go/app/domain/mock"
	"yatter-backend-go/app/domain/repository"
)

// DaoMock is a mock implementation of Dao
type DaoMock struct {
}

func (d *DaoMock) Account() repository.Account {
	return &mock.AccountMock{}
}

func (d *DaoMock) Status() repository.Status {
	return &mock.StatusMock{}
}

func (d *DaoMock) Attachment() repository.Attachment {
	return &mock.AttachmentMock{}
}

func (d *DaoMock) InitAll() error {
	return nil
}
