package request

import (
	"io"
	"net/http"
	"regexp"
)

const (
	Image   string = "image"
	Gifv    string = "gifv"
	Video   string = "video"
	Unknown string = "unknown"
)

func DetectAttachmentType(file io.ReadSeeker) (string, error) {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return classifyFiletype(filetype), nil
}

var imageType = regexp.MustCompile("image/*")
var videoType = regexp.MustCompile("video/*")

func classifyFiletype(filetype string) string {
	if filetype == "image/gif" {
		return Gifv
	} else if imageType.MatchString(filetype) {
		return Image
	} else if videoType.MatchString(filetype) {
		return Video
	}
	return Unknown
}
