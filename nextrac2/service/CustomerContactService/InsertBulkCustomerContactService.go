package CustomerContactService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/MasterDataService/DistrictService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/PersonProfileService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/ProvinceService"
	"strings"
	"time"
)

func (input customerContactService) InsertBulkCustomerContactFromCustomer(tx *sql.Tx, inputStruct []in.CustomerContactRequest, contextModel *applicationModel.ContextModel, timeNow time.Time, isForUpdate, isNewCompanyProfile bool) (output out.CustomerErrorResponse, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var custContactOnDB repository.CustomerContactModel

	for i := 0; i < len(inputStruct); i++ {
		var isUpdate bool
		err = input.validateInsertBulkCustomerContactOnDB(&inputStruct[i], contextModel, false)
		if err.Error != nil {
			output.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
			inputStruct[i].IsSuccess = false
			return
		}

		if isForUpdate {
			custContactOnDB, err = dao.CustomerContactDAO.GetCustomerContactForUpdate(tx, repository.CustomerContactModel{
				Nik:        sql.NullString{String: inputStruct[i].Nik},
				CustomerID: sql.NullInt64{Int64: inputStruct[i].CustomerID},
			})
			if err.Error != nil {
				output.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
				inputStruct[i].IsSuccess = false
				return
			}
			if custContactOnDB.ID.Int64 > 0 {
				if custContactOnDB.Nik.String == inputStruct[i].Nik {
					err = errorModel.GenerateDataUsedError(input.FileName, "InsertBulkCustomerContactFromCustomer", constanta.NIK)
					output.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
					inputStruct[i].IsSuccess = false
					return
				}
			}
		}

		if inputStruct[i].MdbPersonProfileID < 1 {
			isUpdate = false
		} else {
			isUpdate = true
		}

		// validate master data
		err = master_data_dao.ValidatePersonProfileToMasterData(input.convertCustContactToMDBRequest(inputStruct[i]), contextModel, isUpdate)
		if err.Error != nil {
			output.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
			inputStruct[i].IsSuccess = false
			return
		}
	}

	err = input.validateNIK(inputStruct)
	if err.Error != nil {
		return
	}

	for i := 0; i < len(inputStruct); i++ {
		err = input.handleCustomerContactOnMasterData(&inputStruct[i], contextModel, isNewCompanyProfile)
		if err.Error != nil {
			output.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
			//todo error detail master data
			return
		}
	}

	inputModels := input.convertDTOToModelInsertBulk(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	insertedIDs, err := dao.CustomerContactDAO.InsertBulkCustomerContact(tx, inputModels)
	if err.Error != nil {
		return
	}

	for _, id := range insertedIDs {
		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.CustomerContactDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: id},
		})
	}

	return
}

func (input customerContactService) validateInsertBulkCustomerContactOnDB(item *in.CustomerContactRequest, contextModel *applicationModel.ContextModel, isUpdate bool) (err errorModel.ErrorModel) {
	funcName := "validateInsertBulkCustomerContactOnDB"
	var (
		dataPositionOnMDB master_data_response.PositionResponse
		db                = serverconfig.ServerAttribute.DBConnection
		provinceOnDB      repository.ProvinceModel
		districtOnDB      repository.DistrictModel
	)
	//var dataInterface interface{}

	// Get data scope
	scope, err := input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	if !isUpdate {
		// validate Person Profile On MDB
		if item.MdbPersonProfileID > 0 {
			_, err = PersonProfileService.PersonProfileServie.DoViewPersonProfile(
				master_data_request.PersonProfileGetListRequest{
					ID: item.MdbPersonProfileID,
				}, contextModel)
			if err.Error != nil {
				if err.Error.Error() == constanta.ErrorMDBDataNotFound {
					err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.PersonProfile)
				}
				return
			}
		}
	}

	// Validation Person Title
	dataPersonTitleMDB, err := master_data_dao.ViewDetailPersonTitleFromMasterData(int(item.MdbPersonTitleID), contextModel)
	if err.Error != nil {
		service.LogMessage("Error Person Title", http.StatusBadRequest)
		if err.Error.Error() == constanta.ErrorMDBDataNotFound {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.PersonTitle)
		}
		return
	}

	if dataPersonTitleMDB.ID < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.PersonTitle)
		return
	}

	item.PersonTitle = dataPersonTitleMDB.Title

	// Validation Position
	dataPositionOnMDB, err = master_data_dao.GetViewPositionFromMasterData(master_data_request.PositionGetListRequest{
		ID: item.MdbPositionID,
	}, contextModel)
	if err.Error != nil {
		if err.Error.Error() == constanta.ErrorMDBDataNotFound {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Position)
		}
		return
	}

	item.PositionName = dataPositionOnMDB.Name
	//if item.MdbPositionID > 0 {
	//} else {
	//	item.PositionName = ""
	//}

	// Validate Province
	if item.ProvinceID > 0 {
		provinceOnDB, err = dao.ProvinceDAO.GetProvinceForCustomer(db, repository.ProvinceModel{
			ID: sql.NullInt64{Int64: item.ProvinceID},
		}, scope, ProvinceService.ProvinceService.MappingScopeDB)
		if err.Error != nil {
			return
		}
		if provinceOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Province)
			return
		}
		item.MDBProvinceID = provinceOnDB.MDBProvinceID.Int64
	}

	// Validate District ID
	if item.DistrictID > 0 {
		districtOnDB, err = dao.DistrictDAO.GetDistrictWithProvinceID(db, repository.ListLocalDistrictModel{
			ID:         sql.NullInt64{Int64: item.DistrictID},
			ProvinceID: sql.NullInt64{Int64: item.ProvinceID},
		}, scope, DistrictService.DistrictService.MappingScopeDB)
		if err.Error != nil {
			return
		}
		if districtOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.District)
			return
		}
		item.MDBDistrictID = districtOnDB.MdbDistrictID.Int64
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerContactService) validateNIK(inputStruct []in.CustomerContactRequest) (err errorModel.ErrorModel) {
	for i, item := range inputStruct {
		var newArr []in.CustomerContactRequest
		if item.Action == constanta.ActionDeleteCode {
			continue
		}

		newArr = append(newArr, inputStruct[:i]...)
		newArr = append(newArr, inputStruct[(i+1):]...)

		for _, newItem := range newArr {
			if newItem.Nik == item.Nik && !util.IsStringEmpty(item.Nik) && newItem.Action != constanta.ActionDeleteCode {
				err = errorModel.GenerateDataDuplicateInDTOError(input.FileName, "validateNIK")
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerContactService) convertDTOToModelInsertBulk(inputStruct []in.CustomerContactRequest, authAccessModel model.AuthAccessTokenModel, timeNow time.Time) (result []repository.CustomerContactModel) {
	for _, item := range inputStruct {
		result = append(result, repository.CustomerContactModel{
			CustomerID:         sql.NullInt64{Int64: item.CustomerID},
			MdbPersonProfileID: sql.NullInt64{Int64: item.MdbPersonProfileID},
			Nik:                sql.NullString{String: item.Nik},
			MdbPersonTitle:     sql.NullInt64{Int64: item.MdbPersonTitleID},
			PersonTitle:        sql.NullString{String: item.PersonTitle},
			FirstName:          sql.NullString{String: item.FirstName},
			LastName:           sql.NullString{String: item.LastName},
			Sex:                sql.NullString{String: item.Sex},
			Address:            sql.NullString{String: item.Address},
			Address2:           sql.NullString{String: item.Address2},
			Address3:           sql.NullString{String: item.Address3},
			Hamlet:             sql.NullString{String: item.Hamlet},
			Neighbourhood:      sql.NullString{String: item.Neighbourhood},
			ProvinceID:         sql.NullInt64{Int64: item.ProvinceID},
			DistrictID:         sql.NullInt64{Int64: item.DistrictID},
			Phone:              sql.NullString{String: item.Phone},
			Email:              sql.NullString{String: item.Email},
			MdbPositionID:      sql.NullInt64{Int64: item.MdbPositionID},
			PositionName:       sql.NullString{String: item.PositionName},
			Status:             sql.NullString{String: item.Status},
			CreatedBy:          sql.NullInt64{Int64: authAccessModel.ResourceUserID},
			CreatedAt:          sql.NullTime{Time: timeNow},
			CreatedClient:      sql.NullString{String: authAccessModel.ClientID},
			UpdatedBy:          sql.NullInt64{Int64: authAccessModel.ResourceUserID},
			UpdatedAt:          sql.NullTime{Time: timeNow},
			UpdatedClient:      sql.NullString{String: authAccessModel.ClientID},
		})
	}
	return result
}

func (input customerContactService) ValidateInsertBulkCustomerContact(inputStruct *in.CustomerContactRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel, detailResponses out.CustomerErrorResponse) {
	err = inputStruct.ValidateInsert(false)
	if err.Error != nil {
		detailResponses.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct.Nik)
	}
	return
}

func (input customerContactService) convertBulkCustContactToMDBRequest(inputStruct []in.CustomerContactRequest) (output []master_data_request.PersonProfileWriteRequest) {
	for _, request := range inputStruct {
		output = append(output, input.convertCustContactToMDBRequest(request))
	}
	return
}

func (input customerContactService) convertCustContactToMDBRequest(inputStruct in.CustomerContactRequest) (output master_data_request.PersonProfileWriteRequest) {
	arrPhone := strings.Split(inputStruct.Phone, "-")

	return master_data_request.PersonProfileWriteRequest{
		ID:               inputStruct.MdbPersonProfileID,
		PersonTitleID:    inputStruct.MdbPersonTitleID,
		Nik:              inputStruct.Nik,
		FirstName:        inputStruct.FirstName,
		LastName:         inputStruct.LastName,
		Sex:              inputStruct.Sex,
		Address1:         inputStruct.Address,
		Address2:         inputStruct.Address2,
		Address3:         inputStruct.Address3,
		Hamlet:           inputStruct.Hamlet,
		Neighbourhood:    inputStruct.Neighbourhood,
		ProvinceID:       inputStruct.MDBProvinceID,
		DistrictID:       inputStruct.MDBDistrictID,
		PhoneCountryCode: arrPhone[0],
		Phone:            arrPhone[1],
		Email:            inputStruct.Email,
		Status:           constanta.StatusActive,
		UpdatedAt:        inputStruct.UpdatedAt,
	}
}

func (input customerContactService) getPrevPayload(inputStruct []in.CustomerContactRequest) (output out.PreviousPayload) {
	for _, request := range inputStruct {
		output.CustomerContact = append(output.CustomerContact, out.PreviousCustomerContact{
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

func (input customerContactService) handleCustomerContactOnMasterData(inputStruct *in.CustomerContactRequest, contextModel *applicationModel.ContextModel, isNewCompanyProfile bool) (err errorModel.ErrorModel) {
	var personProfile master_data_response.ViewPersonProfileResponse

	if inputStruct.MdbPersonProfileID == 0 {
		//--- Insert person profile
		inputStruct.MdbPersonProfileID, err = master_data_dao.InsertPersonProfileToMasterData(input.convertCustContactToMDBRequest(*inputStruct), contextModel)
		if err.Error != nil {
			inputStruct.IsSuccess = false
			return
		}

		if util.IsStringEmpty(inputStruct.Nik) {
			personProfile, err = master_data_dao.GetViewPersonProfileFromMasterData(master_data_request.PersonProfileGetListRequest{
				ID: inputStruct.MdbPersonProfileID,
			}, contextModel)
			if err.Error != nil {
				inputStruct.IsSuccess = false
				return
			}

			inputStruct.Nik = personProfile.Nik
		}

		//--- Insert contact_person on master data
		inputStruct.MdbContactPersonID, err = master_data_dao.InsertContactPerson(input.convertCustContactToContactPersonRequest(*inputStruct), contextModel)
		if err.Error != nil {
			inputStruct.IsSuccess = false
			return
		}
	} else {
		//--- Validate person_profile on MDB
		personProfile, err = master_data_dao.GetViewPersonProfileFromMasterData(master_data_request.PersonProfileGetListRequest{ID: inputStruct.MdbPersonProfileID}, contextModel)
		if err.Error != nil {
			return
		}

		personProfileTemp := personProfile
		personProfileTemp.ProvinceID, err = dao.ProvinceDAO.GetProvinceIDByMdbID(serverconfig.ServerAttribute.DBConnection,
			repository.ProvinceModel{MDBProvinceID: sql.NullInt64{Int64: personProfile.ProvinceID}}, false)
		if err.Error != nil {
			return
		}

		personProfileTemp.DistrictID, err = dao.DistrictDAO.GetDistrictIDByMdbID(serverconfig.ServerAttribute.DBConnection,
			repository.DistrictModel{
				MdbDistrictID: sql.NullInt64{Int64: personProfile.DistrictID},
				ProvinceID:    sql.NullInt64{Int64: personProfile.ProvinceID},
			})
		if err.Error != nil {
			return
		}

		err = input.validateChangedData(*inputStruct, personProfileTemp)
		if err.Error != nil {
			return
		}

		inputStruct.MdbPersonTitleID = personProfile.PersonTitleID
		inputStruct.FirstName = personProfile.FirstName
		inputStruct.LastName = personProfile.LastName
		inputStruct.Sex = personProfile.Sex
		inputStruct.Address = personProfile.Address1
		inputStruct.Address2 = personProfile.Address2
		inputStruct.Address3 = personProfile.Address3
		inputStruct.Hamlet = personProfile.Hamlet
		inputStruct.Neighbourhood = personProfile.Neighbourhood
		inputStruct.MDBProvinceID = personProfile.ProvinceID
		inputStruct.MDBDistrictID = personProfile.DistrictID
		inputStruct.UpdatedAt = personProfile.UpdatedAt
		inputStruct.ProvinceID = personProfileTemp.ProvinceID
		inputStruct.DistrictID = personProfileTemp.DistrictID

		//inputStruct.ProvinceID, err = dao.ProvinceDAO.GetProvinceIDByMdbID(serverconfig.ServerAttribute.DBConnection,
		//	repository.ProvinceModel{MDBProvinceID: sql.NullInt64{Int64: personProfile.ProvinceID}}, false)
		//if err.Error != nil {
		//	return
		//}

		//inputStruct.DistrictID, err = dao.DistrictDAO.GetDistrictIDByMdbID(serverconfig.ServerAttribute.DBConnection,
		//	repository.DistrictModel{
		//		MdbDistrictID: sql.NullInt64{Int64: personProfile.DistrictID},
		//		ProvinceID:    sql.NullInt64{Int64: personProfile.ProvinceID},
		//	})
		//if err.Error != nil {
		//	return
		//}

		//--- Update person_profile
		err = master_data_dao.UpdatePersonProfileToMasterData(input.convertCustContactToMDBRequest(*inputStruct), contextModel)
		if err.Error != nil {
			inputStruct.IsSuccess = false
			return
		}

		if isNewCompanyProfile {
			inputStruct.MdbContactPersonID, err = master_data_dao.InsertContactPerson(input.convertCustContactToContactPersonRequest(*inputStruct), contextModel)
			if err.Error != nil {
				inputStruct.IsSuccess = false
				return
			}
		} else {
			var isExist bool

			inputStruct.MdbContactPersonID, isExist, err = master_data_dao.ValidateContactPersonOnMDB(master_data_request.ContactPersonGetListRequest{
				ParentID:  inputStruct.MdbCompanyProfileID,
				Connector: constanta.TableNameMasterDataCompanyProfile,
				NIK:       inputStruct.Nik,
			}, contextModel)
			if err.Error != nil {
				return
			}

			if isExist {
				//--- Update contact_person on master data
				err = master_data_dao.UpdateContactPerson(input.convertCustContactToContactPersonRequest(*inputStruct), contextModel)
				if err.Error != nil {
					inputStruct.IsSuccess = false
					return
				}
			} else {
				//--- Insert contact_person on master data
				inputStruct.MdbContactPersonID, err = master_data_dao.InsertContactPerson(input.convertCustContactToContactPersonRequest(*inputStruct), contextModel)
				if err.Error != nil {
					inputStruct.IsSuccess = false
					return
				}
			}
		}
	}

	inputStruct.IsSuccess = true
	return
}

func (input customerContactService) convertCustContactToContactPersonRequest(inputStruct in.CustomerContactRequest) (output master_data_request.ContactPersonWriteRequest) {
	arrPhone := strings.Split(inputStruct.Phone, "-")

	return master_data_request.ContactPersonWriteRequest{
		ID:              inputStruct.MdbContactPersonID,
		PersonProfileID: inputStruct.MdbPersonProfileID,
		PersonTitleID:   inputStruct.MdbPersonTitleID,
		FirstName:       inputStruct.FirstName,
		LastName:        inputStruct.LastName,
		NIK:             inputStruct.Nik,
		Address1:        inputStruct.Address,
		ParentID:        inputStruct.MdbCompanyProfileID,
		PositionID:      inputStruct.MdbPositionID,
		Connector:       constanta.TableNameMasterDataCompanyProfile,
		Email:           inputStruct.Email,
		PhoneCode:       arrPhone[0],
		Phone:           arrPhone[1],
		SuperiorID:      0,
		Status:          constanta.StatusActive,
	}
}
