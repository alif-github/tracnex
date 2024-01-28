package CustomerService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/MasterDataService/CompanyTitleService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/DistrictService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/ProvinceService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
)

type customerService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var CustomerService = customerService{}.New()

func (input customerService) New() (output customerService) {
	output.FileName = "CustomerService.go"
	output.ServiceName = constanta.Customer
	output.ValidLimit = []int{5}
	output.ValidLimit = append(output.ValidLimit, service.DefaultLimit...)
	output.ValidOrderBy = []string{
		"id",
		"customer_name",
		"address",
		"district_name",
		"province_name",
		"phone",
		"status",
	}
	output.ValidSearchBy = []string{
		"id",
		"province_id",
		"district_id",
		"customer_name",
		"customer_group_id",
		"customer_category_id",
		"salesman_id",
		"distributor_of",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.CustomerCategoryDataScope] = applicationModel.MappingScopeDB{
		View:  "c.customer_category_id",
		Count: "c.customer_category_id",
	}
	output.MappingScopeDB[constanta.CustomerGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "c.customer_group_id",
		Count: "c.customer_group_id",
	}
	output.MappingScopeDB[constanta.SalesmanDataScope] = applicationModel.MappingScopeDB{
		View:  "c.salesman_id",
		Count: "c.salesman_id",
	}
	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View:  "c.province_id",
		Count: "c.province_id",
	}
	output.MappingScopeDB[constanta.DistrictDataScope] = applicationModel.MappingScopeDB{
		View:  "c.district_id",
		Count: "c.district_id",
	}

	output.ListScope = input.SetListScope()

	return
}

func (input customerService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.CustomerRequest) errorModel.ErrorModel) (inputStruct in.CustomerRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input customerService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_customer_npwp") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.NPWP)
		}
	}

	return err
}

func (input customerService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{
		constanta.CustomerGroupDataScope,
		constanta.CustomerCategoryDataScope,
		constanta.SalesmanDataScope,
		constanta.ProvinceDataScope,
		constanta.DistrictDataScope,
	})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerService) validateMDBComponent(inputStruct *in.CustomerRequest, scope map[string]interface{}, contextModel *applicationModel.ContextModel, internalToken string) (err errorModel.ErrorModel) {
	var (
		funcName         = "validateMDBComponent"
		db               = serverconfig.ServerAttribute.DBConnection
		dataMDB          interface{}
		provinceOnDB     repository.ProvinceModel
		DistrictOnDB     repository.DistrictModel
		SubDistrictOnDB  repository.SubDistrictModel
		UrbanVillageOnDB repository.UrbanVillageModel
		PostalCodeOnDB   repository.PostalCodeModel
		msg              = util2.GenerateConstantaI18n(constanta.MismatchRegionalData, contextModel.AuthAccessTokenModel.Locale, nil)
	)

	//--- Validate company title on MDB
	if inputStruct.MDBCompanyTitleID > 0 {
		if dataMDB, err = CompanyTitleService.CompanyTitleService.DoViewCompanyTitle(
			master_data_request.CompanyTitleRequest{
				ID: inputStruct.MDBCompanyTitleID,
			}, contextModel); err.Error != nil && err.Error.Error() == constanta.ErrorMDBDataNotFound {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CompanyTitle)
			return
		}

		dataCompanyTitle := dataMDB.(master_data_response.CompanyTitleResponse)
		inputStruct.CompanyTitle = dataCompanyTitle.Title
	} else {
		inputStruct.CompanyTitle = ""
	}

	fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate company title on MDB", funcName))

	//--- Validate province on MDB
	provinceOnDB, err = dao.ProvinceDAO.GetProvinceForCustomer(db, repository.ProvinceModel{
		ID: sql.NullInt64{Int64: inputStruct.ProvinceID},
	}, scope, ProvinceService.ProvinceService.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if provinceOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Province)
		return
	}

	fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Province local", funcName))

	//--- View province on MDB
	_, err = master_data_dao.ViewDetailProvinceFromMasterData(int(provinceOnDB.MDBProvinceID.Int64), contextModel, internalToken)
	if err.Error != nil {
		err = input.checkRegionalDataError(err, funcName, contextModel, constanta.Province)
		return
	}

	fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Province MDB", funcName))

	inputStruct.CountryID = provinceOnDB.CountryID.Int64
	inputStruct.MDBProvinceID = provinceOnDB.MDBProvinceID.Int64

	//--- Validate district on MDB
	DistrictOnDB, err = dao.DistrictDAO.GetDistrictWithProvinceID(db, repository.ListLocalDistrictModel{
		ID:         sql.NullInt64{Int64: inputStruct.DistrictID},
		ProvinceID: sql.NullInt64{Int64: inputStruct.ProvinceID},
	}, scope, DistrictService.DistrictService.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if DistrictOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.District)
		return
	}

	fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate District Local", funcName))

	//--- View district on MDB
	_, err = master_data_dao.ViewDetailDistrictFromMasterData(int(DistrictOnDB.MdbDistrictID.Int64), contextModel, internalToken)
	if err.Error != nil {
		err = input.checkRegionalDataError(err, funcName, contextModel, constanta.District)
		return
	}

	inputStruct.MDBDistrictID = DistrictOnDB.MdbDistrictID.Int64

	fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate District MDB", funcName))

	//--- Validate sub district on MDB
	if inputStruct.SubDistrictID > 0 {
		//--- Validate sub district on MDB
		SubDistrictOnDB, err = dao.SubDistrictDAO.GetSubDistrictWithDistrictID(db, repository.SubDistrictModel{
			ID:         sql.NullInt64{Int64: inputStruct.SubDistrictID},
			DistrictID: sql.NullInt64{Int64: inputStruct.DistrictID},
		})

		if err.Error != nil {
			return
		}

		if SubDistrictOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.SubDistrict)
			return
		}

		fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Sub District Local", funcName))

		//--- View sub district on MDB
		_, err = master_data_dao.GetViewSubDistrictFromMasterData(SubDistrictOnDB.MDBSubDistrictID.Int64, contextModel, internalToken)
		if err.Error != nil {
			if err.Error.Error() == constanta.MasterDataUnknownDataErrorCode {
				err = errorModel.GenerateErrorCustomActivationCode(input.FileName, funcName, fmt.Sprintf(`[%s] %s`, constanta.SubDistrict, msg))
			}
			return
		}

		fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Sub District MDB", funcName))

		inputStruct.MDBSubDistrictID = SubDistrictOnDB.MDBSubDistrictID.Int64
	}

	//--- Validate urban village on MDB
	if inputStruct.UrbanVillageID > 0 {
		//--- Validate urban village on MDB
		UrbanVillageOnDB, err = dao.UrbanVillageDAO.GetUrbanVillageWithSubDistrictID(db, repository.UrbanVillageModel{
			ID:            sql.NullInt64{Int64: inputStruct.UrbanVillageID},
			SubDistrictID: sql.NullInt64{Int64: inputStruct.SubDistrictID},
		})

		if err.Error != nil {
			return
		}

		if UrbanVillageOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UrbanVillage)
			return
		}

		fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Urban Village Local", funcName))

		//--- View urban village on MDB
		_, err = master_data_dao.GetViewUrbanVillageFromMasterData(master_data_request.UrbanVillageRequest{ID: UrbanVillageOnDB.MDBUrbanVillageID.Int64}, contextModel)
		if err.Error != nil {
			if err.Error.Error() == constanta.MasterDataUnknownDataErrorCode {
				err = errorModel.GenerateErrorCustomActivationCode(input.FileName, funcName, fmt.Sprintf(`[%s] %s`, constanta.UrbanVillage, msg))
			}
			return
		}

		fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Urban Village MDB", funcName))

		inputStruct.MDBUrbanVillageID = UrbanVillageOnDB.MDBUrbanVillageID.Int64
	}

	//--- Validate postal code on MDB
	if inputStruct.PostalCodeID > 0 {
		//--- Validate postal code on MDB
		PostalCodeOnDB, err = dao.PostalCodeDAO.GetPostalCodeWithUrbanVillageID(db, repository.PostalCodeModel{
			ID:             sql.NullInt64{Int64: inputStruct.PostalCodeID},
			UrbanVillageID: sql.NullInt64{Int64: inputStruct.UrbanVillageID},
		})

		if err.Error != nil {
			return
		}

		if PostalCodeOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.PostalCode)
			return
		}

		fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Postal Code Local", funcName))

		//--- View postal code on MDB
		_, err = master_data_dao.GetViewPostalCodeFromMasterData(master_data_request.PostalCodeRequest{ID: PostalCodeOnDB.MDBPostalCodeID.Int64}, contextModel)
		if err.Error != nil {
			if err.Error.Error() == constanta.MasterDataUnknownDataErrorCode {
				err = errorModel.GenerateErrorCustomActivationCode(input.FileName, funcName, fmt.Sprintf(`[%s] %s`, constanta.PostalCode, msg))
			}
			return
		}

		fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate Postal Code MDB", funcName))

		inputStruct.MDBPostalCodeID = PostalCodeOnDB.MDBPostalCodeID.Int64
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
