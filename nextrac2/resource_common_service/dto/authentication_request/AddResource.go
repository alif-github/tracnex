package authentication_request

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type AddResourceClient struct {
	ClientID				string							`json:"client_id"`
	ClientSecret			string							`json:"client_secret"`
	ResourceID				string							`json:"resource_id"`
	AdditionalInformation	[]model.AdditionalInformation	`json:"additional_information"`
}
