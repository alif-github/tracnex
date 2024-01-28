package DataGroupService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type dataGroupService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var DataGroupService = dataGroupService{}.New()

func (input dataGroupService) New() (output dataGroupService) {
	output.FileName = "DataGroupService.go"
	output.ServiceName = "DATA_GROUP"
	output.ValidLimit = []int{10, 20, 50, 100, 200, 500}
	output.ValidOrderBy = []string{"created_name", "id", "group_id", "created_at", "description"}
	output.ValidSearchBy = []string{"group_id","description"}
	return
}

func (input dataGroupService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.DataGroupRequest) errorModel.ErrorModel) (inputStruct in.DataGroupRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateErrorFormatJSON(input.FileName, "readBodyAndValidate", errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)

	//input.TrimDTO(&inputStruct)

	err = validation(&inputStruct)

	return
}

func (input dataGroupService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_datagroup_groupid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.DataGroup)
		}
	}

	return err
}
