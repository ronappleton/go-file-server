package files

import (
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/segmentio/fasthash/fnv1a"
	"os"
	"strconv"
	"strings"
)

var mimeTypes = []string{
	"image/jpeg",
	"image/bmp",
	"image/gif",
	"image/png",
	"image/vnd",
	"application/octet-stream",
}

type File struct {
	Key      string
	Name     string
	Data     string
	MimeType string
	FilePath string
	UrlPath  string
}

func GetMimeType(path string) (string, error) {
	mimeType, _ := mimetype.DetectFile(path)
	fmt.Println("MimeType detected: " + mimeType.String())
	if !contains(mimeTypes, mimeType.String()) {
		return "", errors.New("unable to obtain mime type")
	}

	return mimeType.String(), nil
}

func CreateFile(path string) File {
	info, _ := os.Stat(path)
	data, _ := os.ReadFile(path)
	mimeType, _ := GetMimeType(path)

	return File{
		Key:      strconv.FormatUint(fnv1a.HashString64(path), 10),
		Name:     info.Name(),
		Data:     string(data),
		MimeType: mimeType,
		FilePath: path,
		UrlPath:  getUrlPath(path),
	}
}

func getUrlPath(path string) string {
	pathParts := strings.Split(path, "/")
	return trimmedPath(pathParts, "images")
}

func trimmedPath(pathParts []string, needle string) string {
	var idx = 0
	for index, str := range pathParts {
		if str == needle {
			idx = index
		}
	}

	return strings.Join(append(pathParts[:0], pathParts[idx+1:]...), "/")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
