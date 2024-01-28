package util

import "net/http"

func CORSOriginHandler(responseWriter *http.ResponseWriter) {
	(*responseWriter).Header().Set("Access-Control-Allow-Origin", "*")
	(*responseWriter).Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization, x-nextoken, CLIENT_ID, RESOURCE_USER_ID")
	(*responseWriter).Header().Set("Access-Control-Allow-Credentials", "true")
	(*responseWriter).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	(*responseWriter).Header().Set("Access-Control-Max-Age", "1209600")
}
