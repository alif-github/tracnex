package out

type ClientRegisterNonOnPremiseResponse struct {
	UniqueID1    string `json:"unique_id_1"`
	UniqueID2    string `json:"unique_id_2"`
	IsError      bool   `json:"is_error"`
	ErrorMessage string `json:"error_message"`
}
