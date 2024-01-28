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

type companyUpdateService struct {
	FileName string
	service.AbstractService
}

var CompanyUpdateService = companyUpdateService{FileName: "CompanyUpdateService.go"}

func (input companyUpdateService) UpdateCompany(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UpdateCompany"
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	companyBody, files, err := readBodyWithFileAndValidate(request, context)

	if err.Error != nil {
		return
	}

	companyBody.ID = id

	payload := companyStruct{
		inputStruct: companyBody,
		files:       files,
	}

	_, err = input.InsertServiceWithAudit(funcName, payload, context, input.doUpdate, input.uploadFileToAzure)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", context)

	return
}

func (input companyUpdateService) doUpdate(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdate"

	temp := body.(companyStruct)

	companyBody := temp.inputStruct
	files := temp.files

	companyRepo := GetCompanyRepository(companyBody, context, now)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if isOnlyHaveOwnAccess {
		companyRepo.CreatedBy.Int64 = userID
	}

	companyOnDB, err := dao.CompanyDAO.GetDetailCompanyForUpdate(tx, companyRepo.ID.Int64)
	if err.Error != nil {
		return
	}

	if companyOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "company")
		return
	}

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

	if urbanVillageOnDb.ID.Int64 == 0{
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


	err = input.validation(companyOnDB, companyBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.EmployeeLevelDAO.TableName, companyRepo.ID.Int64, userID)...)

	err = dao.CompanyDAO.UpdateCompany(tx, companyRepo)
	if err.Error != nil {
		return
	}

	auditData, err = input.uploadPhotoToLocalCDN(tx, files, companyOnDB.ID.Int64, context, now)
	if err.Error != nil {
		return
	}

	return
}

func (input companyUpdateService) validation(companyOnDB repository.CompanyModel, companyBody in.CompanyRequest) (err errorModel.ErrorModel) {
	err = companyBody.ValidateCompany(true)
	if err.Error != nil {
		return
	}
	err = service.OptimisticLock(companyOnDB.UpdatedAt.Time, companyBody.UpdatedAt, input.FileName, "company")
	return
}

func (input companyUpdateService) uploadFileToAzure(data interface{}, contextModel applicationModel.ContextModel) {
	service.UploadListFileToAzure(data, contextModel, dao.FileUploadDAO.UpdateFileUploads)
}

func (input companyUpdateService) uploadPhotoToLocalCDN(tx *sql.Tx, files []in.MultipartFileDTO, parentId int64, contextModel *applicationModel.ContextModel, now time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	if len(files) > 1 {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "uploadPhotoToLocalCDN", "Foto yang diupload hanya boleh 1", "foto", "")
		return
	}

	photo, err := dao.FileUploadDAO.GetFileByParentIDAndCategory(serverconfig.ServerAttribute.DBConnection, parentId, "kms_category")
	if err.Error != nil {
		return
	}

	if len(files) == 1{
		var photoRepo []in.MultipartFileDTO
		var photoEdit in.MultipartFileDTO
		photoEdit = in.MultipartFileDTO{
			Host: photo.Host.String,
			Path: photo.Path.String,
		}

		photoRepo = append(photoRepo, photoEdit)
		service.DeleteFileFromCDN(photoRepo)

		err = dao.FileUploadDAO.DeleteFileUploadByParentID(tx, parentId, dao.CompanyDAO.TableName,  dao.CompanyDAO.TableName)
		if err.Error != nil {
			return
		}
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