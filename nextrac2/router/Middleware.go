package router

import (
	"context"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

func Middleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		util.CORSOriginHandler(&responseWriter)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Header().Set(constanta.RequestCacheControl, "no-cache")
		if request.Method == "OPTIONS" {
			return
		} else {
			var contextModel *applicationModel.ContextModel
			defer func() {
				if r := recover(); r != nil {
					util2.InputLog(errorModel.GenerateRecoverError(), contextModel.LoggerModel)
				}
			}()

			requestID := request.Header.Get(constanta.RequestIDConstanta)
			if requestID == "" {
				requestID = util.GetUUID()
				request.Header.Set(constanta.RequestIDConstanta, requestID)
			}

			request.Header.Set(constanta.RequestCacheControl, "no-cache")

			var contextModels applicationModel.ContextModel

			contextModels.LoggerModel = applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
			contextModels.LoggerModel.RequestID = requestID
			contextModels.LoggerModel.IP = request.Header.Get(constanta.IPAddressConstanta)
			contextModels.LoggerModel.Source = request.Header.Get(constanta.SourceConstanta)
			contextModels.LoggerModel.Class = "[Middleware.go,Middleware]"

			ctx := context.WithValue(request.Context(), constanta.ApplicationContextConstanta, &contextModels)
			request = request.WithContext(ctx)

			timestamp := time.Now()

			nextHandler.ServeHTTP(responseWriter, request)

			contextModel = request.Context().Value(constanta.ApplicationContextConstanta).(*applicationModel.ContextModel)
			contextModel.LoggerModel.Time = int64(time.Since(timestamp).Seconds())
			logMiddleware(contextModel.LoggerModel, request.RequestURI)
		}
	})
}

func logMiddleware(loggerModel applicationModel.LoggerModel, requestURI string) {
	if !strings.Contains(requestURI, "/health") && !(loggerModel.IP == "" && loggerModel.Class == "[Middleware.go,Middleware]") {
		util2.InputLog(errorModel.GenerateNonErrorModel(), loggerModel)
	}
}
