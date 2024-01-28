package UserParameterService

import (
	//"encoding/json"
	//"github.com/gorilla/mux"
	//"net/http"
	//"nexsoft.co.id/nextrac2/constanta"
	//"nexsoft.co.id/nextrac2/dto/in"
	//"nexsoft.co.id/nextrac2/model/applicationModel"
	//"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	//"strconv"
)

type userParameterService struct {
	service.AbstractService
	service.GetListData
}

var UserParameterService = userParameterService{}.New()

func (input userParameterService) New() (output userParameterService) {
	output.FileName = "UserParameterService.go"
	output.ValidLimit = []int{10, 20, 50, 100, 200, 500}
	output.ValidOrderBy = []string{"id", "permission", "user_id"}
	output.ValidSearchBy = []string{"permission", "user_id"}
	return
}

//func (input userParameterService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UserParameterDTOIn) errorModel.ErrorModel) (inputStruct in.UserParameterDTOIn, err errorModel.ErrorModel) {
//	var stringBody string
//
//	stringBody, err = input.ReadBody(request, contextModel)
//	if err.Error != nil {
//		return
//	}
//
//	_ = json.Unmarshal([]byte(stringBody), &inputStruct)
//
//	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
//	inputStruct.UserID = int64(id)
//
//	err = validation(&inputStruct)
//
//	return
//}

//func (input userParameterService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
//	if err.CausedBy != nil {
//		if service.CheckDBError(err, "uq_userparameter_userid") {
//			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.UserID)
//		}
//	}
//
//	return err
//}
