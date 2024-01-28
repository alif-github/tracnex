package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
	"time"
)

type BacklogRequest struct {
	ID               int64       `json:"id"`
	DepartmentName   string      `json:"department_name"`
	Feature          int64       `json:"feature"`
	ReferenceTicket  int64       `json:"reference_ticket"`
	Subject          string      `json:"subject"`
	Layer1           string      `json:"layer_1"`
	Layer2           string      `json:"layer_2"`
	Layer3           string      `json:"layer_3"`
	Layer4           string      `json:"layer_4"`
	Layer5           string      `json:"layer_5"`
	RedmineNumberStr string      `json:"redmine_number"`
	Sprint           string      `json:"sprint"`
	SprintName       string      `json:"sprint_name"`
	PicId            int64       `json:"pic_id"`
	PicName          string      `json:"pic_name"`
	Status           string      `json:"status"`
	Mandays          float64     `json:"mandays"`
	MandaysDone      float64     `json:"mandays_done"`
	FlowChanged      string      `json:"flow_changed"`
	AdditionalData   string      `json:"additional_data"`
	Tracker          string      `json:"tracker"`
	Note             string      `json:"note"`
	Url              string      `json:"url"`
	Page             string      `json:"page"`
	DepartmentCode   string      `json:"department_code"`
	FileUploadId     int64       `json:"file_upload_id"`
	File             FileReqInfo `json:"file"`
	UpdatedAtStr     string      `json:"updated_at"`
	RedmineNumber    int64
	UpdatedAt        time.Time
}

type FileReqInfo struct {
	FileName string `json:"file_name"`
	Type     string `json:"type"`
	Base64   string `json:"base_64"`
}

type ImportBacklogRequest struct {
	DepartmentCode string `json:"department_code"`
}

type MultipleUpdateStatusRequest struct {
	ID     []int64 `json:"id"`
	Status string  `json:"status"`
}

func (input *BacklogRequest) ValidateView() (err errorModel.ErrorModel) {
	var (
		fileName = "BacklogDTO.go"
		funcName = "ValidateView"
	)

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *BacklogRequest) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		fileName = "BacklogDTO.go"
		funcName = "ValidateDelete"
	)

	return input.validationForUpdateAndDelete(fileName, funcName)
}

func (input *MultipleUpdateStatusRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	var (
		fileName = "BacklogDTO.go"
		funcName = "ValidateUpdate"
	)

	if len(input.ID) < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.Status) {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Status)
	}

	//if util.IsStringEmpty(input.UpdatedAtStr) {
	//	return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	//}
	//
	//input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	//if err.Error != nil {
	//	return
	//}

	// Status
	//err = input.validationAllStatusAllowed()
	//if err.Error != nil {
	//	return
	//}

	return errorModel.GenerateNonErrorModel()
}

func (input *BacklogRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *BacklogRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	var (
		fileName          = "BacklogDTO.go"
		funcName          = "ValidateUpdate"
		departmentAllowed = []string{
			constanta.DepartmentQAQC,
			constanta.DepartmentDeveloper,
		}
	)

	if err = input.validationForUpdateAndDelete(fileName, funcName); err.Error != nil {
		return
	}

	err = input.checkDepartmentAllowed(departmentAllowed)
	if err.Error != nil {
		return
	}

	return input.mappingValidationByDepartment(fileName)
}

func (input *BacklogRequest) ValidateInsert() (err errorModel.ErrorModel) {
	var (
		fileName = "BacklogDTO.go"
	)

	err = input.checkDepartmentAllowed(constanta.ListDepartmentNexsoft)
	if err.Error != nil {
		return
	}

	return input.mappingValidationByDepartment(fileName)
}

func (input *BacklogRequest) checkDepartmentAllowed(departmentAllowed []string) (err errorModel.ErrorModel) {
	var (
		counter = 0
	)

	for _, department := range departmentAllowed {
		if input.DepartmentCode == department {
			counter++
			return
		}
	}

	if counter < 1 {
		err = errorModel.GenerateSimpleErrorModel(400, "Code department tidak diperbolehkan")
	}

	return
}

func (input *BacklogRequest) mappingValidationByDepartment(fileName string) (err errorModel.ErrorModel) {
	var funcName = "mappingValidationByDepartment"

	if input.DepartmentCode == constanta.DepartmentQAQC {
		return input.validationInsertQAQC(fileName, funcName)
	}

	return input.validationInsertDeveloper(fileName, funcName)
}

func (input *BacklogRequest) validationInsertDeveloper(fileName string, funcName string) (err errorModel.ErrorModel) {
	input.Tracker = "Task"
	err = input.globalMandatoryFieldValidation(fileName, funcName)

	if util.IsStringEmpty(input.SprintName) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.SprintName)
	}

	if err.Error != nil {
		return
	}

	return input.optionalFieldInsertDeveloper(fileName, funcName)
}

func (input *BacklogRequest) validationInsertQAQC(fileName string, funcName string) (err errorModel.ErrorModel) {
	err = input.MandatoryFieldInsertQAQCValidation(fileName, funcName)
	if err.Error != nil {
		return
	}

	return input.optionalFieldInsertQAQC(fileName, funcName)
}

func (input *BacklogRequest) globalMandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(input.DepartmentCode) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.DepartmentCode)
	}

	if util.IsStringEmpty(input.Layer1) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Layer1)
	}

	if util.IsStringEmpty(input.Sprint) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Sprint)
	}

	if input.PicId < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.PicId)
	}

	if input.Mandays <= 0 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Mandays)
	}

	return
}

func (input *BacklogRequest) MandatoryFieldInsertQAQCValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	err = input.globalMandatoryFieldValidation(fileName, funcName)
	if err.Error != nil {
		return
	}

	if input.Feature < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Feature)
	}

	err = util2.ValidateMinMaxInteger(input.Feature, constanta.Feature, 1, 10)
	if err.Error != nil {
		return
	}

	//-- Redmine Number
	if util.IsStringEmpty(input.RedmineNumberStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.RedmineNumber)
	}

	var (
		numTemp int
		errorS  error
	)

	r := strings.Trim(input.RedmineNumberStr, "#")
	numTemp, errorS = strconv.Atoi(r)
	if errorS != nil {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "redmine harus angka dengan simbol # di depan, ex: #202308", constanta.RedmineNumber, "")
	}

	input.RedmineNumber = int64(numTemp)
	if input.RedmineNumber < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.RedmineNumber)
	}

	err = util2.ValidateMinMaxInteger(input.Feature, constanta.RedmineNumber, 1, 10)
	if err.Error != nil {
		return
	}

	if util.IsStringEmpty(input.Status) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Status)
	}

	if util.IsStringEmpty(input.Subject) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Subject)
	}

	err = util2.ValidateMinMaxString(input.Subject, constanta.Subject, 1, 256)
	if err.Error != nil {
		return
	}

	if input.Mandays <= 0 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Mandays)
	}

	err = util2.ValidateMinMaxFloat(input.Mandays, constanta.Mandays, 1, 3)
	if err.Error != nil {
		return
	}

	if util.IsStringEmpty(input.Tracker) {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Tracker)
	}

	err = input.validationTrackerAllowed()
	if err.Error != nil {
		return
	}

	return
}

func (input *BacklogRequest) globalOptionalFieldInsert(fileName string, funcName string) (err errorModel.ErrorModel) {

	// Layer 1
	if !util.IsStringEmpty(input.Layer1) {
		err = util2.ValidateMinMaxString(input.Layer1, constanta.Layer1, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Layer1, input.Layer1)
		if err.Error != nil {
			return
		}
	}

	// Layer 2
	if !util.IsStringEmpty(input.Layer2) {
		err = util2.ValidateMinMaxString(input.Layer2, constanta.Layer2, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Layer2, input.Layer2)
		if err.Error != nil {
			return
		}
	}

	// Layer 3
	if !util.IsStringEmpty(input.Layer3) {
		err = util2.ValidateMinMaxString(input.Layer3, constanta.Layer3, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Layer3, input.Layer3)
		if err.Error != nil {
			return
		}
	}

	// Layer 4
	if !util.IsStringEmpty(input.Layer4) {
		err = util2.ValidateMinMaxString(input.Layer4, constanta.Layer4, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Layer4, input.Layer4)
		if err.Error != nil {
			return
		}
	}

	// Layer 5
	if !util.IsStringEmpty(input.Layer5) {
		err = util2.ValidateMinMaxString(input.Layer5, constanta.Layer5, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Layer5, input.Layer5)
		if err.Error != nil {
			return
		}
	}

	// Redmine
	if !util.IsStringEmpty(input.RedmineNumberStr) {
		var (
			numTemp int
			errorS  error
		)

		r := strings.Trim(input.RedmineNumberStr, "#")
		numTemp, errorS = strconv.Atoi(r)
		if errorS != nil {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "redmine harus angka dengan simbol # di depan, ex: #202308", constanta.RedmineNumber, "")
		}

		input.RedmineNumber = int64(numTemp)
		if input.RedmineNumber < 1 {
			return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.RedmineNumber)
		}

		err = util2.ValidateMinMaxInteger(input.RedmineNumber, constanta.RedmineNumber, 1, 10)
		if err.Error != nil {
			return
		}
	}

	// Sprint Name
	if !util.IsStringEmpty(input.SprintName) {
		err = util2.ValidateMinMax(input.SprintName, constanta.SprintName, 1, 30)
		if err.Error != nil {
			return
		}

		isValid, errMsg, addInfo := IsOnlyAlfaNumerikValid(input.SprintName)
		if !isValid {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errMsg, constanta.SprintName, addInfo)
			return
		}
	}

	// Status
	//if !util.IsStringEmpty(input.Status) {
	//	err = input.validationStatusAllowed()
	//	if err.Error != nil {
	//		return
	//	}
	//}

	// Mandays
	if input.Mandays > 0 {
		err = util2.ValidateMinMaxFloat(input.Mandays, constanta.Mandays, 1, 3)
		if err.Error != nil {
			return
		}
	}

	// Mandays Done
	if input.MandaysDone > 0 {
		err = util2.ValidateMinMaxFloat(input.MandaysDone, constanta.MandaysDone, 1, 3)
		if err.Error != nil {
			return
		}
	}

	// Flow Changed
	if !util.IsStringEmpty(input.FlowChanged) {
		err = util2.ValidateMinMaxString(input.FlowChanged, constanta.FlowChanged, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.FlowChanged, input.FlowChanged)
		if err.Error != nil {
			return
		}
	}

	// Additional Data
	if !util.IsStringEmpty(input.AdditionalData) {
		err = util2.ValidateMinMaxString(input.AdditionalData, constanta.AdditionalData, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.AdditionalData, input.AdditionalData)
		if err.Error != nil {
			return
		}
	}

	// Note
	if !util.IsStringEmpty(input.Note) {
		err = util2.ValidateMinMaxString(input.Note, constanta.Note, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Note, input.Note)
		if err.Error != nil {
			return
		}
	}

	// Url
	if !util.IsStringEmpty(input.Url) {
		err = util2.ValidateMinMaxString(input.Url, constanta.Url, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Url, input.Url)
		if err.Error != nil {
			return
		}
	}

	// Page
	if !util.IsStringEmpty(input.Page) {
		err = util2.ValidateMinMaxString(input.Page, constanta.Page, 1, 256)
		if err.Error != nil {
			return
		}

		err = util2.ValidateEmptyField(fileName, funcName, constanta.Page, input.Page)
		if err.Error != nil {
			return
		}
	}

	// Sprint
	if !util.IsStringEmpty(input.Sprint) {
		if len(input.Sprint) != 17 {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "FIXED", constanta.Sprint, strconv.Itoa(17)+" Character.")
		}
	}

	return
}

func (input *BacklogRequest) validationStatusAllowed() (err errorModel.ErrorModel) {

	if input.DepartmentCode == constanta.DepartmentQAQC {
		return input.validateStatus(constanta.StatusAllowedQA)
	}

	return input.validateStatus(constanta.StatusAllowedDeveloper)
}

func (input *MultipleUpdateStatusRequest) validationAllStatusAllowed() (err errorModel.ErrorModel) {
	var allStatus = append(constanta.StatusAllowedDeveloper)
	allStatus = append(allStatus, constanta.StatusAllowedQA...)

	return input.validateStatus(allStatus)
}

func (input *BacklogRequest) validationTrackerAllowed() (err errorModel.ErrorModel) {
	var (
		counter = 0
	)

	for _, tracker := range constanta.ListTrackerQA {
		if input.Tracker == tracker {
			counter++
		}
	}

	if counter < 1 {
		err = errorModel.GenerateSimpleErrorModel(400, "Tracker tidak diperbolehkan")
	}

	return
}

func (input *BacklogRequest) validateStatus(statusAllowed []string) (err errorModel.ErrorModel) {
	var counter = 0

	for _, status := range statusAllowed {
		if input.Status == status {
			counter++
		}
	}

	if counter < 1 {
		err = errorModel.GenerateSimpleErrorModel(400, "Status tidak diperbolehkan")
	}

	return
}

func (input *MultipleUpdateStatusRequest) validateStatus(statusAllowed []string) (err errorModel.ErrorModel) {
	var counter = 0

	for _, status := range statusAllowed {
		if input.Status == status {
			counter++
		}
	}

	if counter < 1 {
		err = errorModel.GenerateSimpleErrorModel(400, "Status tidak diperbolehkan")
	}

	return
}

func (input *BacklogRequest) optionalFieldInsertQAQC(fileName string, funcName string) (err errorModel.ErrorModel) {
	err = input.globalOptionalFieldInsert(fileName, funcName)
	if err.Error != nil {
		return
	}

	// Reference Ticket
	if input.ReferenceTicket > 0 {
		err = util2.ValidateMinMaxInteger(input.ReferenceTicket, constanta.ReferenceTicket, 1, 10)
		if err.Error != nil {
			return
		}
	}

	// Subject
	if !util.IsStringEmpty(input.Subject) {
		err = util2.ValidateMinMaxString(input.Subject, constanta.Subject, 1, 256)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *BacklogRequest) optionalFieldInsertDeveloper(fileName string, funcName string) (err errorModel.ErrorModel) {
	return input.globalOptionalFieldInsert(fileName, funcName)
}
