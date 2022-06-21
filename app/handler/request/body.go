package request

import (
	"io"
	"net/http"
	"regexp"
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
		return "gifv"
	} else if imageType.MatchString(filetype) {
		return "image"
	} else if videoType.MatchString(filetype) {
		return "video"
	}
	return "unknown"
}
