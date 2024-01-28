package util

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/OneOfOne/xxhash"
	"strconv"
)

func CheckSumWithXXHASH(content []byte) (checksum string) {
	hash := xxhash.Checksum64(content)
	return strconv.Itoa(int(hash))
}

func CheckSumWithMD5(content []byte) (checksum string) {
	hash := md5.New()
	hash.Write(content)
	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes)
}

func CheckSumWithSha256(content []byte) string {
	result := sha256.Sum256(content)
	return hex.EncodeToString(result[:])
}

func CheckSumWithSha512(content []byte) string {
	result := sha512.Sum512(content)
	return hex.EncodeToString(result[:])
}
