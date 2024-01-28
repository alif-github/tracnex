package BacklogService

import (
	"bytes"
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/mux"
	"mime/multipart"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type backlogService struct {
	service.AbstractService
	service.GetListData
}

var BacklogService = backlogService{}.New()

func (input backlogService) New() (output backlogService) {
	output.FileName = "BacklogService.go"
	output.ServiceName = "BACKLOG"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"sprint"}
	output.ValidSearchBy = []string{"sprint"}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "b.employee_id",
		Count: "b.employee_id",
	}
	return
}

func (input backlogService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.BacklogRequest) errorModel.ErrorModel) (inputStruct in.BacklogRequest, err errorModel.ErrorModel) {
	var (
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)

	return
}

func (input backlogService) readBodyAndValidateForMultipleUpdateStatusBacklog(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.MultipleUpdateStatusRequest) errorModel.ErrorModel) (inputStruct in.MultipleUpdateStatusRequest, err errorModel.ErrorModel) {
	var (
		fileName   = "BacklogService.go"
		funcName   = "readBodyAndValidateForMultipleUpdateStatusBacklog"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errs := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errs != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errs)
		return
	}

	err = validation(&inputStruct)

	return
}

func (input backlogService) readBodyAndValidateForInsert(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.BacklogRequest) errorModel.ErrorModel) (inputStruct []in.BacklogRequest, err errorModel.ErrorModel) {
	var (
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	for _, inputStructItem := range inputStruct {
		err = validation(&inputStructItem)
	}

	return
}

func (input backlogService) readBodyAndValidateFileCSVBacklog(request *http.Request, inputStruct in.ImportBacklogRequest, contextModel *applicationModel.ContextModel) (records *excelize.File, inputs in.ImportBacklogRequest, err errorModel.ErrorModel) {
	var (
		errs error
		file multipart.File
	)

	// Get File From Request
	file, _, errs = request.FormFile("file-backlog")
	if errs != nil {
		err = errorModel.GenerateSimpleErrorModel(400, "File tidak dapat dibaca")
		return
	}

	if file == nil {
		err = errorModel.GenerateSimpleErrorModel(400, "Harap lampirkan file backlog")
		return
	}
	defer file.Close()

	content := request.MultipartForm.Value["content"][0]
	byteContent := []byte(content)
	_ = json.Unmarshal(byteContent, &inputStruct)

	//------- Validation input request content
	inputs, err = input.validationDepartmentCode(inputStruct)
	if err.Error != nil {
		return
	}

	// Read Excel File
	records, errs = excelize.OpenReader(file)
	if errs != nil {
		err = errorModel.GenerateSimpleErrorModel(500, "File tidak sesuai format")
		return
	}

	return
}

func (input backlogService) validationDepartmentCode(inputStruct in.ImportBacklogRequest) (inputs in.ImportBacklogRequest, err errorModel.ErrorModel) {
	var departmentsAllowed = []string{constanta.DepartmentQAQC, constanta.DepartmentDeveloper}
	var counter int64

	for _, departmentAllowed := range departmentsAllowed {
		if inputStruct.DepartmentCode == departmentAllowed {
			counter++
		}
	}

	if counter < 1 {
		err = errorModel.GenerateSimpleErrorModel(400, "Department code tidak diperbolehkan")
	}

	inputs = inputStruct
	return
}

func (input backlogService) ReadRequestMultipartForm(request *http.Request, contextModel *applicationModel.ContextModel,
	validation func(input []*in.BacklogRequest) errorModel.ErrorModel) (inputStruct []*in.BacklogRequest, buffer *bytes.Buffer, extension string, err errorModel.ErrorModel) {

	var (
		fileName = "InsertBacklogDetailService.go"
		funcName = "readRequestMultipartForm"
		errs     error
		//file     multipart.File
		//handler  *multipart.FileHeader
	)

	//------- Maximum 10 MB
	errs = request.ParseMultipartForm(constanta.SizeMaximumRequireFileImport)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	//------- Get request from multipart form
	//file, handler, errorS = request.FormFile(config.ApplicationConfiguration.GetDataDirectory().KeyFile)
	//if errorS != nil {
	//	err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
	//	return
	//}

	//------- Get request content
	content := request.MultipartForm.Value["content"][0]
	byteContent := []byte(content)
	errs = json.Unmarshal(byteContent, &inputStruct)
	if errs != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errs)
		return
	}

	//------- Validation input request content
	//for _, inputStructItem := range inputStruct {
	//	//
	//	//}
	//	//
	//	//for i:=0; i < len(inputStruct)-1 ; i++ {
	//	//	err = validation(&inputStruct, inputStruct[i], i)
	//	//	if err.Error != nil {
	//	//		return
	//	//	}
	//	//}

	err = validation(inputStruct)
	if err.Error != nil {
		return
	}

	//defer func() {
	//	_ = file.Close()
	//}()

	//------- Write to contextModel
	//contextModel.LoggerModel.ByteIn = int(handler.Size)

	//------- Get extension with split
	//extensionTemp := strings.Split(handler.Filename, ".")
	//if len(extensionTemp) == 1 {
	//	err = errorModel.GenerateFormatFieldError(fileName, funcName, "file")
	//	return
	//}

	//------- Error if extension is not "csv"
	//if extensionTemp[len(extensionTemp)-1] != "csv" {
	//	err = errorModel.GenerateFormatFieldError(fileName, funcName, "file")
	//	return
	//}

	//extension = extensionTemp[len(extensionTemp)-1]

	//------- Buffer read
	//buffer = bytes.NewBuffer(nil)
	//if _, errorS = io.Copy(buffer, file); errorS != nil {
	//	err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
	//	return
	//}

	return
}

func (input backlogService) validateDataScope(contextModel *applicationModel.ContextModel) (scope map[string]interface{}, err errorModel.ErrorModel) {
	return input.ValidateMultipleDataScope(contextModel, []string{constanta.EmployeeDataScope})
}

func (input backlogService) ReadAndValidateDeleteListDataBacklog(request *http.Request, validSearchKey []string, validOperator map[string]applicationModel.DefaultOperator, validLimit []int) (inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	inputStruct = input.readDeleteListData(request)

	searchByParam, err = inputStruct.ValidateDeleteListData(validSearchKey, validOperator, validLimit)
	return
}

func (input backlogService) readDeleteListData(request *http.Request) (inputStruct in.GetListDataDTO) {
	inputStruct.Search = service.GenerateQueryValue(request.URL.Query()["filter"])
	return
}

func (input backlogService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_redmine_number_backlog") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.RedmineNumber)
		}
	}

	return err
}
