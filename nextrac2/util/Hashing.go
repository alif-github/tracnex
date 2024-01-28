package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func GetHmacSHA512(input string, key string) string {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(input))
	return hex.EncodeToString(mac.Sum(nil))
}

func GenerateSHA256(input string) string {
	sha256 := sha256.New()
	sha256.Write([]byte(input))
	return hex.EncodeToString(sha256.Sum(nil))
}
