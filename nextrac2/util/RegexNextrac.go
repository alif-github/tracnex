package util

import (
	"regexp"
)

func IsClientIDValid(clientID string) (bool, string, string) {
	parentClientID := regexp.MustCompile("^[0-9a-f]{32}")
	return parentClientID.MatchString(clientID), "CLIENT_ID_REGEX_MESSAGE", ""
}