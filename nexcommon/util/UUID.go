package util

import (
	"github.com/google/uuid"
	"strings"
)

func GetUUID() (output string) {
	UUID, _ := uuid.NewRandom()
	output = UUID.String()
	output = strings.Replace(output, "-", "", -1)
	return
}
