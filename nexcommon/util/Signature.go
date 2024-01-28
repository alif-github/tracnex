package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateMessageDigest(message string) string {
	message = strings.Replace(message,"\r","",-1)
	return CheckSumWithSha256([]byte(message))
}

func GenerateSignature(httpMethod string, relativeURL string, accessToken string, messageDigest string, timestamp string, key string) string {
	message := httpMethod + ":" + relativeURL + ":" + accessToken + ":" + messageDigest + ":" + timestamp
	return HMAC256(message, key)
}

func ValidateSignature(httpMethod string, relativeURL string, accessToken string, messageDigest string, timestamp string, key string, signature string) bool {
	message := httpMethod + ":" + relativeURL + ":" + accessToken + ":" + messageDigest + ":" + timestamp
	return signature == HMAC256(message, key)
}

func HMAC256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
