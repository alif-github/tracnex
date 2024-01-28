package PKCEClientMappingEndpoint

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/endpoint"
)

type pkceClientMappingEndpoints struct {
	FileName string
	endpoint.AbstractEndpoint
}

var PKCEClientMappingEndpoints pkceClientMappingEndpoints

func (input pkceClientMappingEndpoints) getMenuCodePKCEClientMapping() string {
	return endpoint.GetMenuCode(constanta.MenuUserMasterConsumerPKCEClientMappingRedesign, constanta.MenuUserMasterConsumerPKCEClientMapping)
}
