package dao

import (
	"yatter-backend-go/app/domain/mock"
	"yatter-backend-go/app/domain/repository"
)

// DaoMock is a mock implementation of Dao
type DaoMock struct {
	AccountMock    *mock.AccountMock
	StatusMock     *mock.StatusMock
	AttachmentMock *mock.AttachmentMock
}

func NewMock(accountMock *mock.AccountMock, statusMock *mock.StatusMock, attachmentMock *mock.AttachmentMock) *DaoMock {
	return &DaoMock{
		AccountMock:    accountMock,
		StatusMock:     statusMock,
		AttachmentMock: attachmentMock,
	}
}

func (d *DaoMock) Account() repository.Account {
	return d.AccountMock
}

func (d *DaoMock) Status() repository.Status {
	return d.StatusMock
}

func (d *DaoMock) Attachment() repository.Attachment {
	return d.AttachmentMock
}

func (d *DaoMock) InitAll() error {
	return nil
}
