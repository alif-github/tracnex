package util

import (
	"encoding/base64"
	"errors"
)

func Base64decoder(content string) ([]byte, error) {
	output, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		output, err = base64.RawStdEncoding.DecodeString(content)
		if err != nil {
			inside := content
			inside = inside + "=="
			output, err = base64.StdEncoding.DecodeString(inside)
			if err != nil {
				return nil, errors.New("HASHING_DATA_INVALID")
			}
		}
	}
	return output, err
}

func Base64encoder(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}
