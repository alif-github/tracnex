package out

import "nexsoft.co.id/nextrac2/model/errorModel"

type StatFile struct {
	FileSuccess		int						`json:"file_success"`
	FileFailed		int						`json:"file_failed"`
	FileAmount		int						`json:"file_amount"`
	FileDetail		[]ImportFile			`json:"file_detail"`
}

type ImportFile struct {
	NumberData		int						`json:"number_data"`
	ErrorMessage	errorModel.ErrorModel	`json:"error_message"`
}

type MultipleErrorResponse struct {
	ID			int64
	CausedBy 	string
}

type ImportDataResponse struct {
	Filename 	string
	TotalData	int
	Data     	[][]string
}