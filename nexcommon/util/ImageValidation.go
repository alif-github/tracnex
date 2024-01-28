package util

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func IsFileImage(fileContent []byte, maxSize int) (error, string) {
	if len(fileContent) > maxSize {
		return errors.New("NEED_LESS_THAN"), strconv.Itoa(maxSize) + " byte"
	}

	if !mimeFromIncipit(fileContent) {
		return errors.New("PHOTO_REGEX"), ""
	}

	return nil, ""
}

func mimeFromIncipit(incipit []byte) bool {
	var magicTable = map[string]string{
		"\xff\xd8\xff":      "image/jpeg",
		"\x89PNG\r\n\x1a\n": "image/png",
	}
	incipitStr := string(incipit)
	for magic, _ := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return true
		}
	}

	return false
}
