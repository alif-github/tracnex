package util

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func ReadBody(request *http.Request) (output string, bodySize int, err error) {
	byteBody, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		return "", 0, errors.New("BODY_INVALID")
	}
	return string(byteBody), len(byteBody), nil
}

func ReadHeader(request *http.Request, headerName string) (result string) {
	result = request.Header.Get(headerName)
	if result == "" {
		result = request.Header.Get(strings.ToLower(headerName))
	}
	return result
}
