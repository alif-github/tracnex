package WhitelistDeviceService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type whitelistDeviceService struct {
	service.AbstractService
	service.GetListData
}

var WhitelistDeviceService = whitelistDeviceService{}.New()

func (input whitelistDeviceService) New() (output whitelistDeviceService) {
	output.FileName = "WhitelistDeviceService.go"
	output.ServiceName = "WHITELIST_DEVICE"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{"device", "description"}
	output.ValidOrderBy = []string{
		"id",
		"device",
		"description",
		"updated_at",
		"updated_by",
		"created_at",
		"updated_name",
	}

	//output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	//output.MappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{
	//	View:  "pg.id",
	//	Count: "pg.id",
	//}
	//
	//output.ListScope = input.SetListScope()

	return
}

func (input whitelistDeviceService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(whiteListDeviceRequest *in.WhiteListDeviceRequest) errorModel.ErrorModel) (inputStruct in.WhiteListDeviceRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		var errorS = json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateErrorFormatJSON(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)

	err = validation(&inputStruct)
	return
}