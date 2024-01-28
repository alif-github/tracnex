package CompanyService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

type companyStruct struct {
	inputStruct  in.CompanyRequest
	files        []in.MultipartFileDTO
	deletedPhoto []in.MultipartFileDTO
}

func getCompanyBody(request *http.Request, fileName string) (inputStruct in.CompanyRequest, bodySize int, err errorModel.ErrorModel) {
	funcName := "getCompanyBody"
	jsonString, bodySize, readError := util.ReadBody(request)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	readError = json.Unmarshal([]byte(jsonString), &inputStruct)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func readBodyWithFileAndValidate(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.CompanyRequest, files []in.MultipartFileDTO, err errorModel.ErrorModel) {
	var totalSize int64

	errs := request.ParseMultipartForm(32 << 20)
	if errs != nil {
		err = errorModel.GenerateUnknownError("CompanyUtils", "readBodyWithFileAndValidate", errs)
		return
	}

	content := request.FormValue("content")
	byteContent := []byte(content)
	_ = json.Unmarshal(byteContent, &inputStruct)

	//files, totalSize, err = service.ReadFileWithMultipart(request, 0, imageValidation)
	files, totalSize, err = service.ReadFileWithMultipart(request, constanta.InternalCompanyAttachmentMaximumLogo, imageValidation)
	if err.Error != nil {
		return
	}
	contextModel.LoggerModel.ByteIn = int(totalSize) + len(byteContent)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	//err = validation(&inputStruct)

	return
}

func imageValidation(dto in.MultipartFileDTO) errorModel.ErrorModel {
	errs, additional := util.IsFileImage(dto.FileContent, constanta.CompanyProfileMaximumPhotoSize)
	if errs != nil {
		return errorModel.GenerateFieldFormatWithRuleError("CompanyUtilsService", "imageValidation", errs.Error(), "", additional)
		//return errorModel.GenerateFieldFormatWithRuleError("CompanyUtilsService", "imageValidation", errs.Error(), constanta.InternalCompanyAttachment, additional)
	}
	return errorModel.GenerateNonErrorModel()
}

func GetCompanyRepository(body in.CompanyRequest, tokenModel *applicationModel.ContextModel, now time.Time) repository.CompanyModel {
	return repository.CompanyModel{
		ID:            sql.NullInt64{Int64:   body.ID},
		CompanyTitle:  sql.NullString{String: body.CompanyTitle},
		CompanyName:   sql.NullString{String: body.CompanyName},
		Address:       sql.NullString{String: body.Address},
		Address2:      sql.NullString{String: body.Address2},
		Neighbourhood: sql.NullString{String: body.Neighbourhood},
		Hamlet:        sql.NullString{String: body.Hamlet},
		ProvinceID:    sql.NullInt64{Int64: body.ProvinceId},
		DistrictID:    sql.NullInt64{Int64: body.DistrictId},
		SubDistrictID: sql.NullInt64{Int64: body.SubDistrictId},
		UrbanVillageID:sql.NullInt64{Int64: body.VillageId},
		PostalCodeID:  sql.NullInt64{Int64: body.PostalCodeId},
		Longitude:     sql.NullString{String: body.Longitude},
		Latitude:      sql.NullString{String: body.Latitude},
		Telephone:          sql.NullString{String: body.Telephone},
		AlternateTelephone: sql.NullString{String: body.TelephoneAlternate},
		Fax:                sql.NullString{String: body.Fax},
		CompanyEmail:       sql.NullString{String: body.Email},
		AlternativeCompanyEmail: sql.NullString{String: body.AlternateEmail},
		Npwp:          sql.NullString{String: body.Npwp},
		TaxName:       sql.NullString{String: body.TaxName},
		TaxAddress:    sql.NullString{String: body.TaxAddress},
		UpdatedClient: sql.NullString{String: tokenModel.AuthAccessTokenModel.ClientID},
		CreatedClient: sql.NullString{String: tokenModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: tokenModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: tokenModel.AuthAccessTokenModel.ResourceUserID},
	}
}