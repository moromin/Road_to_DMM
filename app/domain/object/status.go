package object

// Status status
type Status struct {
	// The internal ID of the status
	ID int64 `json:"id"`

	// The internal ID of the account who posted status
	AccountID int `json:"-" db:"account_id"`

	// The account of posting status
	Account Account `json:"account"`

	// The contents of status
	Content string `json:"content,omitempty"`

	// The time the status was created
	CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`

	// The attachment of status
	// media_attachments Attachment
	MediaAttachments []Attachment `json:"media_attachments,omitempty"`
}
