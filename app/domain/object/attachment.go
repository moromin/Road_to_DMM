package object

type (
	// Attachment attachment
	Attachment struct {
		// The internal ID of attachment
		ID int64 `json:"id"`

		// The type of attachment
		// One of: "image", "video", "gifv", "unknown"
		Type string `json:"type"`

		// The URL of image
		URL string `json:"url"`

		// The description of the image
		Description string `json:"description"`
	}
)
