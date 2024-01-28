package out

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/model/applicationModel"
)

type APIResponse struct {
	Nexsoft NexsoftMessage `json:"nexsoft"`
}

type NexsoftMessage struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

func (ar APIResponse) String() string {
	return util.StructToJSON(ar)
}

type Payload struct {
	Status StatusResponse `json:"status"`
	Data   PayloadData    `json:"data"`
	Other  interface{}    `json:"other"`
}

type StatusResponse struct {
	Success bool     `json:"success"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Detail  []string `json:"detail"`
}

type PayloadData struct {
	Meta    interface{} `json:"meta"`
	Content interface{} `json:"content"`
}

type Header struct {
	RequestID string `json:"request_id"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type SearchByParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type InitiateGetListDataDTOOut struct {
	ValidOrderBy  []string                                    `json:"valid_order_by"`
	ValidSearchBy []string                                    `json:"valid_search_by"`
	ValidSearchParam []SearchByParam                          `json:"valid_search_by_params"`
	ValidLimit    []int                                       `json:"valid_limit"`
	ValidOperator map[string]applicationModel.DefaultOperator `json:"valid_operator"`
	EnumData      interface{}                                 `json:"enum_data"`
	CountData     int                                         `json:"count_data"`
}

type PhotoList struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Host     string `json:"host"`
}
