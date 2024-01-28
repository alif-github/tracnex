package CustomerContactService

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input customerContactService) UpdateCustomerContactForCustomer(tx *sql.Tx, inputStruct []in.CustomerContactRequest, contextModel *applicationModel.ContextModel, timeNow time.Time, customerID, companyProfileID int64) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		dataInsert      []in.CustomerContactRequest
		dataUpdate      []in.CustomerContactRequest
		dataDelete      []in.CustomerContactRequest
		tempAudit       []repository.AuditSystemModel
		detailResponses out.CustomerErrorResponse
		funcName        = "UpdateCustomerContactForCustomer"
	)

	err = input.validateNIK(inputStruct)
	if err.Error != nil {
		return
	}

	for i := 0; i < len(inputStruct); i++ {
		inputStruct[i].CustomerID = customerID
		inputStruct[i].MdbCompanyProfileID = companyProfileID

		switch inputStruct[i].Action {
		case constanta.ActionInsertCode:
			dataInsert = append(dataInsert, inputStruct[i])
			break
		case constanta.ActionUpdateCode:
			dataUpdate = append(dataUpdate, inputStruct[i])
			break
		case constanta.ActionDeleteCode:
			dataDelete = append(dataDelete, inputStruct[i])
			break
		case 0:
			break
		default:
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ActionCode)
			detailResponses.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
			output = detailResponses
			return
		}
	}

	// todo Delete
	if len(dataDelete) > 0 {
		output, tempAudit, err = input.DoMultiDeleteCustomerContact(tx, dataDelete, contextModel, timeNow, true)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, tempAudit...)
	}

	// todo Insert
	if len(dataInsert) > 0 {
		output, tempAudit, err = input.insertCustomerContactForUpdateCustomer(tx, dataInsert, contextModel, timeNow, true)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, tempAudit...)
	}

	// todo Update
	if len(dataUpdate) > 0 {
		output, tempAudit, err = input.doUpdateCustomerContact(tx, dataUpdate, contextModel, timeNow)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, tempAudit...)
	}
	return
}

func (input customerContactService) doUpdateCustomerContact(tx *sql.Tx, inputStruct []in.CustomerContactRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		tempAudit       []repository.AuditSystemModel
		detailResponses out.CustomerErrorResponse
	)

	for i := 0; i < len(inputStruct); i++ {

		output, tempAudit, err = input.updateCustomerContact(tx, inputStruct[i], contextModel, timeNow)
		if err.Error != nil {
			detailResponses.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
			output = detailResponses
			return
		}

		dataAudit = append(dataAudit, tempAudit...)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerContactService) updateCustomerContact(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName           = "updateCustomerContact"
		inputStruct        = inputStructInterface.(in.CustomerContactRequest)
		personProfileOnMDB []master_data_response.GetListPersonProfileResponse
		custContactOnDB    repository.CustomerContactModel
		isNewPersonProfile bool
	)

	err = input.validateUpdateCustomerContact(&inputStruct)
	if err.Error != nil {
		return
	}

	custContactOnDB, err = dao.CustomerContactDAO.GetCustomerContactForUpdate(tx, repository.CustomerContactModel{
		ID:         sql.NullInt64{Int64: inputStruct.ID},
		CustomerID: sql.NullInt64{Int64: inputStruct.CustomerID},
	})
	if err.Error != nil {
		return
	}

	if custContactOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CustomerContact)
		return
	}

	//if custContactOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
	//	err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Customer)
	//	return
	//}

	if custContactOnDB.Nik.String != inputStruct.Nik {
		personProfileOnMDB, err = master_data_dao.GetListPersonProfileFromMasterData(master_data_request.PersonProfileGetListRequest{
			NIK:    inputStruct.Nik,
			Status: constanta.StatusActive,
			AbstractDTO: in.AbstractDTO{
				Page:    1,
				Limit:   10,
				OrderBy: "",
			},
		}, contextModel)

		if err.Error != nil {
			return
		}

		if personProfileOnMDB != nil {
			err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.NIK)
			return
		}

		isNewPersonProfile = true
	}

	//--- validate customer contact
	err = input.validateInsertBulkCustomerContactOnDB(&inputStruct, contextModel, true)
	if err.Error != nil {
		return
	}

	err = input.handleUpdateCompanyProfile(&inputStruct, contextModel, isNewPersonProfile)
	if err.Error != nil {
		return
	}

	inputModel := input.convertDTOToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
	err = dao.CustomerContactDAO.UpdateCustomerContact(tx, inputModel)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerContactDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	return
}

func (input customerContactService) convertDTOToModelUpdate(inputStruct in.CustomerContactRequest, authAccessModel model.AuthAccessTokenModel, timeNow time.Time) (result repository.CustomerContactModel) {
	return repository.CustomerContactModel{
		ID:             sql.NullInt64{Int64: inputStruct.ID},
		CustomerID:     sql.NullInt64{Int64: inputStruct.CustomerID},
		Nik:            sql.NullString{String: inputStruct.Nik},
		MdbPersonTitle: sql.NullInt64{Int64: inputStruct.MdbPersonTitleID},
		PersonTitle:    sql.NullString{String: inputStruct.PersonTitle},
		FirstName:      sql.NullString{String: inputStruct.FirstName},
		LastName:       sql.NullString{String: inputStruct.LastName},
		Sex:            sql.NullString{String: inputStruct.Sex},
		Address:        sql.NullString{String: inputStruct.Address},
		Hamlet:         sql.NullString{String: inputStruct.Hamlet},
		Neighbourhood:  sql.NullString{String: inputStruct.Neighbourhood},
		ProvinceID:     sql.NullInt64{Int64: inputStruct.ProvinceID},
		DistrictID:     sql.NullInt64{Int64: inputStruct.DistrictID},
		Phone:          sql.NullString{String: inputStruct.Phone},
		Email:          sql.NullString{String: inputStruct.Email},
		MdbPositionID:  sql.NullInt64{Int64: inputStruct.MdbPositionID},
		PositionName:   sql.NullString{String: inputStruct.PositionName},
		Status:         sql.NullString{String: inputStruct.Status},
		UpdatedBy:      sql.NullInt64{Int64: authAccessModel.ResourceUserID},
		UpdatedAt:      sql.NullTime{Time: timeNow},
		UpdatedClient:  sql.NullString{String: authAccessModel.ClientID},
	}
}

func (input customerContactService) validateUpdateCustomerContact(inputStruct *in.CustomerContactRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}

func (input customerContactService) getDataAction(dataOnDB []repository.CustomerContactModel, dataInput []in.CustomerContactRequest) (dataInsert []in.CustomerContactRequest, dataUpdate []in.CustomerContactRequest, dataDelete []repository.CustomerContactModel) {
	for _, itemOnDB := range dataOnDB {
		for j, itemInput := range dataInput {
			if itemInput.Nik == itemOnDB.Nik.String {
				dataUpdate = append(dataUpdate, itemInput)
				dataInput = append(dataInput[:j], dataInput[(j+1):]...)
				break
			}
		}
	}

	for _, itemUpdated := range dataUpdate {
		for j, itemOnDB := range dataOnDB {
			if itemOnDB.Nik.String == itemUpdated.Nik {
				dataOnDB = append(dataOnDB[:j], dataOnDB[(j+1):]...)
			}
		}
	}

	dataDelete = dataOnDB
	dataInsert = dataInput
	return
}

func (input customerContactService) handleUpdateCompanyProfile(inputStruct *in.CustomerContactRequest, contextModel *applicationModel.ContextModel, isNewPersonProfile bool) (err errorModel.ErrorModel) {
	//--- Validate person_profile Update
	err = master_data_dao.ValidatePersonProfileToMasterData(input.convertCustContactToMDBRequest(*inputStruct), contextModel, !isNewPersonProfile)
	if err.Error != nil {
		return
	}

	err = input.handleCustomerContactOnMasterData(inputStruct, contextModel, false)
	return
}

func (input customerContactService) validateChangedData(inputStruct in.CustomerContactRequest, customerContactOnDB master_data_response.ViewPersonProfileResponse) (err errorModel.ErrorModel) {
	funcName := "validateChangedData"

	byteInputStruct, _ := json.Marshal(inputStruct)
	byteCusContactDB, _ := json.Marshal(customerContactOnDB)
	fmt.Println("Input Struct : ", string(byteInputStruct))
	fmt.Println("Customer Contact On DB : ", string(byteCusContactDB))

	if customerContactOnDB.PersonTitleID != inputStruct.MdbPersonTitleID {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.MDBPersonTitleID)
		return
	}

	if customerContactOnDB.FirstName != inputStruct.FirstName {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.FirstName)

		return
	}

	if customerContactOnDB.LastName != inputStruct.LastName {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.LastName)
		return
	}

	if customerContactOnDB.Sex != inputStruct.Sex {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.Sex)
		return
	}

	if customerContactOnDB.Address1 != inputStruct.Address {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.Address)
		return
	}

	if customerContactOnDB.Hamlet != inputStruct.Hamlet {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.Hamlet)
		return
	}

	if customerContactOnDB.Neighbourhood != inputStruct.Neighbourhood {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.Neighbourhood)
		return
	}

	if customerContactOnDB.ProvinceID != inputStruct.ProvinceID {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.Province)
		return
	}

	if customerContactOnDB.DistrictID != inputStruct.DistrictID {
		err = errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.District)
		return
	}

	return
}
