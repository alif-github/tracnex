package in

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type NotifyAddClientDTOIn struct {
	ClientID string                      `json:"client_id"`
	RoleName string                      `json:"role_name"`
	Others   model.AdditionalInformation `json:"others"`
}
