package in

import (
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type ClientRegisterLogRequest struct {
	ID                    int64  `json:"id"`
	ClientID              string `json:"client_id"`
	ClientTypeID          int64  `json:"client_type_id"`
	AttributeRequest	  string `json:"attribute_request"`
	SuccessStatusAuth     bool   `json:"success_status_auth"`
	SuccessStatusNexcloud bool   `json:"success_status_nexcloud"`
	SuccessStatusNexdrive bool   `json:"success_status_nexdrive"`
	Resource              string `json:"resource"`
	MessageAuth           string `json:"message_auth"`
	MessageNexcloud       string `json:"message_nexcloud"`
	MessageNexdrive       string `json:"message_nexdrive"`
	Details               string `json:"details"`
	Code                  string `json:"code"`
	RequestTimestamp      time.Time
	RequestCount		  int64	 `json:"request_count"`
}

type ResultWorkerClientLogDTO struct {
	ID						int64
	ClientID				string
	Resource				string
	Errors					errorModel.ErrorModel
}

type PreparedRepositoryClientRegisterLog struct {
	ID                 int64
	ClientID           string
	Status             bool
	Resource           string
	ProcessForResource string
	Code               string
	Message            string
	Detail             string
	AttributeRequest   string
	UpdatedAt          time.Time
}

type PreparedDaoUpdateClientRegisterLog struct {
	StatusKey			string
	MessageKey			string
	StatusValue			bool
	MessageValue		string
}