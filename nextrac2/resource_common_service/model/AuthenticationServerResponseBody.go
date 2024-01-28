package model

type HeaderResponse struct {
	RequestID string `json:"request_id"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type PayloadResponse struct {
	Status PayloadStatusResponse `json:"status"`
	Other  interface{}           `json:"other"`
}

type PayloadStatusResponse struct {
	Success bool     `json:"success"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Detail  []string `json:"info"`
}
