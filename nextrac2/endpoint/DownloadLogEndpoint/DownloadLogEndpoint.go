package DownloadLogEndpoint

import (
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/service/LogFileService"
)

type downloadLogEndpoint struct {
	endpoint.AbstractEndpoint
}

var DownloadLogEndpoint downloadLogEndpoint

func (input downloadLogEndpoint) DownloadLogEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	funcName := "DownloadLogEndpoint.go"
	input.ServeWhiteListEndpointWithFile(funcName, false, responseWriter, request, LogFileService.DownloadLogService.StartService)
}
