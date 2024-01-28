package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type DetailClientMappingAuthenticationResponse struct {
	Nexsoft DetailClientMappingBodyResponse `json:"nexsoft"`
}

type DetailClientMappingBodyResponse struct {
	Header  model.HeaderResponse 			`json:"header"`
	Payload RequestChangePasswordPayload  	`json:"payload"`
}

type DetailClientMappingPayload struct {
	model.PayloadResponse
	Data DetailClientMappingData 			`json:"data"`
}

type DetailClientMappingData struct {
	Content DetailClientMappingContent 		`json:"content"`
}

type DetailClientMappingContent struct {
	ClientID        string                  `json:"client_id"`
	ClientTypeID    int64                   `json:"client_type_id"`
	AuthUserId      int64                   `json:"auth_user_id"`
	UserName        string                  `json:"username"`
	CompanyId       string                  `json:"company_id"`
	BranchId        string                  `json:"branch_id"`
	Aliases         string                  `json:"aliases"`
}