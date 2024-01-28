package CompanyService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type companyInsertService struct {
	FileName string
	service.AbstractService
}

var CompanyInsertService = companyInsertService{FileName: "CompanyInsertService.go"}

func (input companyInsertService) InsertCompany(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	companyBody, files, err := readBodyWithFileAndValidate(request, context)

	if err.Error != nil {
		return
	}

	err = companyBody.ValidateCompany(false)
	if err.Error != nil {
		return
	}

	payload := companyStruct{
		inputStruct: companyBody,
		files:       files,
	}

	_, err = input.InsertServiceWithAudit("InsertCompany", payload, context, input.doInsert, input.uploadFileToAzure)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", context)

	return
}

func (input companyInsertService) doInsert(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {

	temp := body.(companyStruct)

	companyBody := temp.inputStruct
	files := temp.files
	companyRepo := GetCompanyRepository(companyBody, context, now)

	provinceDb, err := dao.ProvinceDAO.CheckIsProvinceExist(serverconfig.ServerAttribute.DBConnection, repository.ProvinceModel{
		ID: sql.NullInt64{Int64: companyRepo.ProvinceID.Int64},
	})

	if err.Error != nil {
		return
	}

	if provinceDb == 0{
		err = errorModel.GenerateUnknownDataError(input.FileName, "doInsert", "province")
		return
	}

	districtOnDb, err := dao.DistrictDAO.GetDistrictByID(serverconfig.ServerAttribute.DBConnection, repository.DistrictModel{
		ID: sql.NullInt64{Int64: companyRepo.DistrictID.Int64},
	})

	if err.Error != nil {
		return
	}

	if districtOnDb.ID.Int64 == 0{
		err = errorModel.GenerateUnknownDataError(input.FileName, "doInsert", "district")
		return
	}

	subDistrictOnDb, err := dao.SubDistrictDAO.GetSubDistrictByIDForUpdate(serverconfig.ServerAttribute.DBConnection, repository.SubDistrictModel{
		ID: sql.NullInt64{Int64: companyRepo.SubDistrictID.Int64},
	})

	if err.Error != nil {
		return
	}

	if subDistrictOnDb.ID.Int64 == 0{
		err = errorModel.GenerateUnknownDataError(input.FileName, "doInsert", "sub_district")
		return
	}

	urbanVillageOnDb, err := dao.UrbanVillageDAO.GetUrbanVillageByIDForGetList(serverconfig.ServerAttribute.DBConnection, companyRepo.UrbanVillageID.Int64, false)

	if err.Error != nil {
		return
	}

	if urbanVillageOnDb.ID.Int64 == 00{
		err = errorModel.GenerateUnknownDataError(input.FileName, "doInsert", "urban_village")
		return
	}

	postalCodeOnDb, err := dao.PostalCodeDAO.GetPostalCodeByID(serverconfig.ServerAttribute.DBConnection, repository.PostalCodeModel{
		ID: sql.NullInt64{Int64: companyRepo.PostalCodeID.Int64},
	})

	if err.Error != nil {
		return
	}

	if postalCodeOnDb.ID.Int64 == 0 && companyRepo.PostalCodeID.Int64 != 0{
		err = errorModel.GenerateUnknownDataError(input.FileName, "doInsert", "postal_code")
		return
	}

	companyId, errs := dao.CompanyDAO.InsertCompany(tx, companyRepo)
	if errs.Error != nil {
		err = errs
		return
	}

	auditData = append(auditData, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.CompanyDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: companyId},
	})

	auditData, err = input.uploadPhotoToLocalCDN(tx, files, companyId, context, now)
	if err.Error != nil {
		return
	}

	return
}

func (input companyInsertService) uploadPhotoToLocalCDN(tx *sql.Tx, files []in.MultipartFileDTO, parentId int64, contextModel *applicationModel.ContextModel, now time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	if len(files) > 1 {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "uploadPhotoToLocalCDN", "Foto yang diupload hanya boleh 1", "foto", "")
		return
	}

	container := constanta.InternalCompanyAttachmentPrefix + service.GetAzureDateContainer()
	err = service.UploadFileToLocalCDN(container, &files, contextModel.AuthAccessTokenModel.ResourceUserID)
	if err.Error != nil {
		return
	}

	for i := 0; i < len(files); i++ {
		var photoID int64
		photo := repository.FileUpload{
			ParentID:      sql.NullInt64{Int64: parentId},
			Category:      sql.NullString{String: "internal_company"},
			Konektor:     sql.NullString{String:  "internal_company"},
			FileName:      sql.NullString{String: files[i].Filename},
			Path:          sql.NullString{String: files[i].Path},
			Host:          sql.NullString{String: files[i].Host},
			CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			CreatedAt:     sql.NullTime{Time: now},
		}

		photoID, err = dao.FileUploadDAO.InsertFileUploadInfoForBacklog(tx, photo)
		if err.Error != nil {
			return
		}

		files[i].FileID = photoID
		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.FileUploadDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: photoID},
		})
	}
	return
}

func (input companyInsertService) uploadFileToAzure(data interface{}, contextModel applicationModel.ContextModel) {
	service.UploadListFileToAzure(data, contextModel, dao.FileUploadDAO.UpdateFileUploads)
}