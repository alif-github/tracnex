package ProcessFileListCustomerService

import (
	"database/sql"
	"errors"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
	"strings"
	"time"
)

func (input processFileListCustomerService) SaveListCustomerInDB(rowData []string, index int, contextModel *applicationModel.ContextModel,
	timeNow time.Time, tx *sql.Tx) (failedIndex int, err errorModel.ErrorModel){

	var mappingData repository.CustomerListModel

	//------- Do validate basic like validate DTO
	mappingData, err = input.validateBasicRowDataListCustomer(rowData[0])
	if err.Error != nil {
		failedIndex = index
		return
	}

	//------- Insert Data
	err = input.insertDataListCustomerToDB(mappingData, contextModel, timeNow, tx)
	if err.Error != nil {
		failedIndex = index
		return
	}

	failedIndex = 0
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) insertDataListCustomerToDB(customerModel repository.CustomerListModel, contextModel *applicationModel.ContextModel,
	timeNow time.Time, tx *sql.Tx) (err errorModel.ErrorModel) {

	fileName := "InsertUpdateListCustomerService.go"
	funcName := "insertDataListCustomerToDB"

	//------- Repo model
	customerModel.CreatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
	customerModel.CreatedAt.Time = timeNow
	customerModel.CreatedClient.String = contextModel.AuthAccessTokenModel.ClientID
	customerModel.UpdatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
	customerModel.UpdatedAt.Time = timeNow
	customerModel.UpdatedClient.String = contextModel.AuthAccessTokenModel.ClientID

	//------- Insert data customer list to DB
	var id int64
	id, err = dao.CustomerListDAO.InsertCustomerByImport(tx, customerModel)
	if err.Error != nil {
		return
	}

	if id < 1 {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errors.New("data tidak ter-insert"))
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) updateDataListCustomerToDB(mapping map[string]interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {

	//------- Repo model
	customerModel := repository.CustomerListModel{
		CompanyID: 		sql.NullString{String: mapping[constanta.MapCompanyID].(string)},
		BranchID: 		sql.NullString{String: mapping[constanta.MapBranchID].(string)},
		CompanyName: 	sql.NullString{String: mapping[constanta.MapCompanyName].(string)},
		City: 			sql.NullString{String: mapping[constanta.MapCity].(string)},
		Implementer: 	sql.NullString{String: mapping[constanta.MapImplementer].(string)},
		Implementation: sql.NullTime{Time: input.dateConvert(mapping[constanta.MapImplementation].(string))},
		Product: 		sql.NullString{String: mapping[constanta.MapProduct].(string)},
		Version:		sql.NullString{String: mapping[constanta.MapVersion].(string)},
		LicenseType: 	sql.NullString{String: mapping[constanta.MapLicenseType].(string)},
		UserAmount: 	sql.NullInt64{Int64: mapping[constanta.MapUserAmount].(int64)},
		ExpDate: 		sql.NullTime{Time: input.dateConvert(mapping[constanta.MapExpDate].(string))},
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt: 		sql.NullTime{Time: timeNow},
		UpdatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	//------- Update data customer list to DB
	err = dao.CustomerListDAO.UpdateCustomer(serverconfig.ServerAttribute.DBConnection, customerModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) validateBasicRowDataListCustomer(rowDataStr string) (customerModel repository.CustomerListModel, err errorModel.ErrorModel) {
	fileName := "InsertUpdateListCustomerService.go"
	funcName := "validateBasicRowDataListCustomer"

	var validationResult bool
	var timeResultImplementation time.Time
	var timeResultExpDate time.Time
	var parseIntegerStr int

	//------- Separate the row data
	dataColTemp := strings.Split(rowDataStr, "|")

	//------- Validate company id
	validationResult = util.IsStringEmpty(dataColTemp[0])
	if validationResult {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyID)
		return
	}

	err = input.ValidateMinMaxString(dataColTemp[0], constanta.CompanyID, 1, 20)
	if err.Error != nil {
		return
	}

	//------- Validate branch id
	validationResult = util.IsStringEmpty(dataColTemp[1])
	if validationResult {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.BranchID)
		return
	}

	err = input.ValidateMinMaxString(dataColTemp[1], constanta.BranchID, 1, 20)
	if err.Error != nil {
		return
	}

	//------- Validate company name
	validationResult = util.IsStringEmpty(dataColTemp[2])
	if validationResult {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyName)
		return
	}

	err = input.ValidateMinMaxString(dataColTemp[2], constanta.CompanyName, 1, 255)
	if err.Error != nil {
		return
	}

	//------- Validate city
	if dataColTemp[3] != "" {
		err = input.ValidateMinMaxString(dataColTemp[3], constanta.City, 1, 50)
		if err.Error != nil {
			return
		}
	}

	//------- Validate implementer
	if dataColTemp[4] != "" {
		err = input.ValidateMinMaxString(dataColTemp[4], constanta.Implementer, 1, 50)
		if err.Error != nil {
			return
		}
	}

	//------- Validate implementation
	if dataColTemp[5] != "" {
		timeResultImplementation = input.dateConvert(dataColTemp[5])
		validationResult = timeResultImplementation.IsZero()

		if validationResult {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.ImplementationAt)
			return
		}
	}

	//------- Validate product
	validationResult = util.IsStringEmpty(dataColTemp[6])
	if validationResult {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Product)
		return
	}

	err = input.ValidateMinMaxString(dataColTemp[6], constanta.Product, 1, 50)
	if err.Error != nil {
		return
	}

	//------- Validate version
	if dataColTemp[7] != "" {
		err = input.ValidateMinMaxString(dataColTemp[7], constanta.Version, 1, 50)
		if err.Error != nil {
			return
		}
	}

	//------- Validate license type
	if dataColTemp[8] != "" {
		err = input.ValidateMinMaxString(dataColTemp[8], constanta.LicenseType, 1, 50)
		if err.Error != nil {
			return
		}
	}

	//------- Validate user amount
	parseIntegerStr, _ = strconv.Atoi(dataColTemp[9])

	//------- Validate exp date
	var expDateFillColumn int
	for i := 17; i >= 10; i-- {
		if dataColTemp[i] != "" {
			expDateFillColumn = i
			break
		}
	}

	timeResultExpDate = input.dateConvert(dataColTemp[expDateFillColumn])
	validationResult = timeResultExpDate.IsZero()
	if validationResult {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.ExpDate)
		return
	}

	customerModel = repository.CustomerListModel{
		CompanyID: 		sql.NullString{String: dataColTemp[0]},
		BranchID: 		sql.NullString{String: dataColTemp[1]},
		CompanyName: 	sql.NullString{String: dataColTemp[2]},
		City: 			sql.NullString{String: dataColTemp[3]},
		Implementer: 	sql.NullString{String: dataColTemp[4]},
		Implementation: sql.NullTime{Time: timeResultImplementation},
		Product: 		sql.NullString{String: dataColTemp[6]},
		Version:		sql.NullString{String: dataColTemp[7]},
		LicenseType: 	sql.NullString{String: dataColTemp[8]},
		UserAmount: 	sql.NullInt64{Int64: int64(parseIntegerStr)},
		ExpDate: 		sql.NullTime{Time: timeResultExpDate},
	}

	err = errorModel.GenerateNonErrorModel()
	return
}