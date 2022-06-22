package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

// Implementation for repository.Status
type status struct {
	db *sqlx.DB
}

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// CreateAccount : content, accountIDから新しいステータスを作成
func (r *status) Create(ctx context.Context, accountID int64, content string, attachmentIDs []int64) (int64, error) {
	var id int64

	err := Transaction(r.db, func(tx *sqlx.Tx) error {
		const registerStatus = `INSERT INTO status (account_id, content) VALUES (?, ?)`
		res, err := tx.ExecContext(ctx, registerStatus, accountID, content)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}

		if len(attachmentIDs) == 0 {
			return nil
		}

		count := 0
		findAttachments, params, err := sqlx.In(`SELECT COUNT(*) FROM attachment WHERE id IN (?)`, attachmentIDs)
		if err != nil {
			return err
		}
		if err := tx.QueryRowxContext(ctx, findAttachments, params...).Scan(&count); err != nil {
			return err
		} else if count != len(attachmentIDs) {
			return fmt.Errorf("attachments %v specified by 'media_ids' is invalid", attachmentIDs)
		}

		type StatusAttachment struct {
			StatusID     int64 `db:"status_id"`
			AttachmentID int64 `db:"attachment_id"`
		}
		statusAttachments := make([]StatusAttachment, len(attachmentIDs))
		for i, attachmentID := range attachmentIDs {
			statusAttachments[i] = StatusAttachment{
				StatusID:     id,
				AttachmentID: attachmentID,
			}
		}
		const registerStatusAttachment = `INSERT INTO status_attachment (status_id, attachment_id) VALUES (:status_id, :attachment_id)`
		if _, err := tx.NamedExecContext(ctx, registerStatusAttachment, statusAttachments); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

// FindByID : IDからステータスを取得
func (r *status) FindByID(ctx context.Context, id int64) (*object.Status, error) {
	status := &object.Status{}

	const findStatus = `SELECT * FROM status WHERE id = ?`
	err := r.db.QueryRowxContext(ctx, findStatus, id).StructScan(status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	attachments, err := r.findAttachments(ctx, id)
	if err != nil {
		return nil, err
	}
	status.MediaAttachments = attachments

	return status, nil
}

func (r *status) findAttachments(ctx context.Context, id int64) ([]object.Attachment, error) {
	const query = `SELECT a.*
					FROM status_attachment as sa
					JOIN attachment as a
					ON sa.attachment_id = a.id
					WHERE sa.status_id = ?`
	attachments := []object.Attachment{}
	if err := r.db.SelectContext(ctx, &attachments, query, id); err != nil {
		return nil, err
	}
	return attachments, nil
}

// DeleteByID : IDからステータスを削除
func (r *status) DeleteByID(ctx context.Context, id int64) error {
	query := `delete
			  from status
			  where id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("status %d is not found", id)
	}

	return nil
}

// ListAll : maxID, sinceID, limit からタイムライン（ステータスのスライス）を取得
func (r *status) ListAll(ctx context.Context, maxID, sinceID, limit int64) ([]object.Status, error) {
	idRange, _ := BuildRangeQuery("s.id", maxID, sinceID, 0)
	listAll := fmt.Sprintf(`SELECT s.*, a.username AS "account.username", a.followers_count AS "account.followers_count", a.following_count AS "account.following_count", a.create_at AS "account.create_at"
							FROM status as s
							JOIN account as a
							on s.account_id = a.id
							%s
							ORDER BY s.id
							LIMIT %d`, idRange, limit)
	statuses := []object.Status{}
	if err := r.db.SelectContext(ctx, &statuses, listAll); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	for i, status := range statuses {
		attachments, err := r.findAttachments(ctx, status.ID)
		if err != nil {
			return nil, err
		}
		statuses[i].MediaAttachments = attachments
	}

	return statuses, nil
}

// ListByID : 認証されたアカウントのID, maxID, sinceID, limit からタイムライン（ステータスのスライス）を取得
func (r *status) ListByID(ctx context.Context, id, maxID, sinceID, limit int64) ([]object.Status, error) {
	connection := ""
	idRange, ok := BuildRangeQuery("id", maxID, sinceID, 0)
	if ok {
		connection = "AND"
	} else {
		connection = "WHERE"
	}
	listByID := fmt.Sprintf(`SELECT * FROM status %s %s account_id = %d LIMIT %d`, idRange, connection, id, limit)
	statuses := []object.Status{}
	if err := r.db.SelectContext(ctx, &statuses, listByID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	for i, status := range statuses {
		attachments, err := r.findAttachments(ctx, status.ID)
		if err != nil {
			return nil, err
		}
		statuses[i].MediaAttachments = attachments
	}

	return statuses, nil
}
