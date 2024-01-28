package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type ElasticSyncRequest struct {
	All            bool `json:"all"`
	Bank           bool `json:"bank"`
	CompanyProfile bool `json:"company_profile"`
	CompanyTitle   bool `json:"company_title"`
	Country        bool `json:"country"`
	District       bool `json:"district"`
	Island         bool `json:"island"`
	PersonProfile  bool `json:"person_profile"`
	PersonTitle    bool `json:"person_title"`
	PostalCode     bool `json:"postal_code"`
	Province       bool `json:"province"`
	SubDistrict    bool `json:"sub_district"`
	UrbanVillage   bool `json:"urban_village"`
}

func (input *ElasticSyncRequest) ValidateElasticSyncRequest() errorModel.ErrorModel {
	var count = 0
	if !input.All {
		if input.Bank {
			count++
		}
		if input.CompanyProfile {
			count++
		}
		if input.CompanyTitle {
			count++
		}
		if input.Country {
			count++
		}
		if input.District {
			count++
		}
		if input.Island {
			count++
		}
		if input.PersonProfile {
			count++
		}
		if input.PersonTitle {
			count++
		}
		if input.PostalCode {
			count++
		}
		if input.Province {
			count++
		}
		if input.SubDistrict {
			count++
		}
		if input.UrbanVillage {
			count++
		}
	} else {
		count++
	}

	if count == 0 {
		return errorModel.GenerateEmptyFieldError("ElasticSyncDTO.go", "ValidateElasticSyncRequest", constanta.ElasticKey)
	}

	return errorModel.GenerateNonErrorModel()
}
