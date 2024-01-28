package util

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func GenerateCryptoRandom() string {
	var (
		maxValueLen = 19
		result      string
	)

	for true {
		number, err := rand.Int(rand.Reader, big.NewInt(9223372036854775807))
		if err == nil {
			result = strconv.Itoa(int(number.Int64()))
			break
		}
	}

	if len(result) < maxValueLen {
		addition := ""
		for i := 0; i < maxValueLen-len(result); i++ {
			addition += "0"
		}
		result = addition + result
	}

	return result
}
