package ReportAnnualLeaveService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input reportAnnualLeaveService) uploadAttachmentFile(tx *sql.Tx, files []in.MultipartFileDTO, contextModel *applicationModel.ContextModel, timeNow time.Time) (fileUploadID int64, err errorModel.ErrorModel) {
	var container = constanta.ContainerReportEmployeeLeave + service.GetAzureDateContainer()
	if len(files) < 1 {
		return
	}

	//--- Upload To CDN
	if err = service.UploadFileToLocalCDN(container, &files, contextModel.AuthAccessTokenModel.ResourceUserID); err.Error != nil {
		return
	}

	//--- Upload To Azure
	go service.UploadFileToAzure(&files)
	fileUploadModel := repository.FileUpload{
		FileName:      sql.NullString{String: files[0].Filename},
		FileSize:      sql.NullInt64{Int64: files[0].Size},
		Konektor:      sql.NullString{String: dao.EmployeeLeaveDAO.TableName},
		Host:          sql.NullString{String: files[0].Host},
		Path:          sql.NullString{String: files[0].Path},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	return dao.FileUploadDAO.InsertFileUploadInfoForBacklog(tx, fileUploadModel)
}
