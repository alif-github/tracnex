package CustomerService

import (
	"errors"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

func getErrorMessage(err errorModel.ErrorModel, contextModel *applicationModel.ContextModel, identityCode string) (output out.CustomerErrorStatus) {
	errCode := err.Error.Error()
	errMessage := util2.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale)
	if errMessage == errCode {
		if err.CausedBy != nil {
			errMessage = err.CausedBy.Error()
		}
	}
	return out.CustomerErrorStatus{
		Npwp:    identityCode,
		Status:  constanta.JobProcessErrorStatus,
		Message: errMessage,
	}
}

func (input customerService) getPreviousPayload(inputStruct in.CustomerRequest) (output out.PreviousPayload) {
	output = out.PreviousPayload{
		ID:                          inputStruct.ID,
		IsPrincipal:                 inputStruct.IsPrincipal,
		IsParent:                    inputStruct.IsParent,
		ParentCustomerID:            inputStruct.ParentCustomerID,
		MDBParentCustomerID:         inputStruct.MDBParentCustomerID,
		MDBCompanyProfileID:         inputStruct.MDBCompanyProfileID,
		Npwp:                        inputStruct.Npwp,
		MDBCompanyTitleID:           inputStruct.MDBCompanyTitleID,
		CompanyTitle:                inputStruct.CompanyTitle,
		CustomerName:                inputStruct.CustomerName,
		Address:                     inputStruct.Address,
		Hamlet:                      inputStruct.Hamlet,
		Neighbourhood:               inputStruct.Neighbourhood,
		CountryID:                   inputStruct.CountryID,
		ProvinceID:                  inputStruct.ProvinceID,
		DistrictID:                  inputStruct.DistrictID,
		SubDistrictID:               inputStruct.SubDistrictID,
		UrbanVillageID:              inputStruct.UrbanVillageID,
		PostalCodeID:                inputStruct.PostalCodeID,
		Longitude:                   inputStruct.Longitude,
		Latitude:                    inputStruct.Latitude,
		PhoneCountryCode:            inputStruct.PhoneCountryCode,
		Phone:                       inputStruct.Phone,
		AlternativePhoneCountryCode: inputStruct.AlternativePhoneCountryCode,
		AlternativePhone:            inputStruct.Phone,
		Fax:                         inputStruct.Fax,
		CompanyEmail:                inputStruct.CompanyEmail,
		AlternativeCompanyEmail:     inputStruct.AlternativeCompanyEmail,
		CustomerSource:              inputStruct.CustomerSource,
		TaxName:                     inputStruct.TaxName,
		TaxAddress:                  inputStruct.Address,
		SalesmanID:                  inputStruct.SalesmanID,
		RefCustomerID:               inputStruct.RefCustomerID,
		DistributorOF:               inputStruct.DistributorOF,
		CustomerGroupID:             inputStruct.CustomerGroupID,
		CustomerCategoryID:          inputStruct.CustomerCategoryID,
		Status:                      inputStruct.Status,
		CustomerContact:             input.getPrevContactPayload(inputStruct.CustomerContact),
		UpdatedAt:                   inputStruct.UpdatedAt,
		IsSuccess:                   inputStruct.IsSuccess,
	}

	return
}

func (input customerService) getPrevContactPayload(inputStruct []in.CustomerContactRequest) (output []out.PreviousCustomerContact) {
	for _, request := range inputStruct {
		output = append(output, out.PreviousCustomerContact{
			ID:                 request.ID,
			CustomerID:         request.CustomerID,
			MdbPersonProfileID: request.MdbPersonProfileID,
			Nik:                request.Nik,
			MdbPersonTitleID:   request.MdbPersonTitleID,
			PersonTitle:        request.PersonTitle,
			FirstName:          request.FirstName,
			LastName:           request.LastName,
			Sex:                request.Sex,
			Address:            request.Address,
			Hamlet:             request.Hamlet,
			Neighbourhood:      request.Neighbourhood,
			ProvinceID:         request.ProvinceID,
			DistrictID:         request.DistrictID,
			Phone:              request.Phone,
			Email:              request.Email,
			MdbPositionID:      request.MdbPositionID,
			PositionName:       request.PositionName,
			Status:             request.Status,
			Action:             request.Action,
			UpdatedAt:          request.UpdatedAt,
			IsSuccess:          request.IsSuccess,
		})
	}

	return
}

func (input customerService) checkRegionalDataError(err errorModel.ErrorModel, funcName string, contextModel *applicationModel.ContextModel, fieldName string) errorModel.ErrorModel {
	errNew := errors.New(constanta.MasterDataUnknownDataErrorCode)
	if err.Error.Error() == errNew.Error() {
		msg := util2.GenerateConstantaI18n(constanta.MismatchRegionalData, contextModel.AuthAccessTokenModel.Locale, nil)
		err = errorModel.GenerateErrorCustomActivationCode(input.FileName, funcName, fmt.Sprintf(`[%s] %s`, fieldName, msg))
	}

	return err
}
