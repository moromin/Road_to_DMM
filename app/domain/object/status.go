package object

// Status status
type Status struct {
	// The internal ID of the status
	ID int64 `json:"-"`

	// The account of posting status
	Account Account `json:"account"`

	// The contents of status
	Content string `json:"content,omitempty"`

	// The time the status was created
	CreateAt DateTime `json:"creat_at,omitempty"`

	// The attachment of status
	// media_attachments Attachment
}
