package UserParameterService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input userParameterService) ViewUserParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	//var inputStruct in.UserParameterDTOIn

	//inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	//output.Data.Content, err = input.doViewUserParameter(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

//func (input userParameterService) doViewUserParameter(inputStruct in.UserParameterDTOIn, contextModel *applicationModel.ContextModel) (result out.ViewDetailUserParameterDTOOut, err errorModel.ErrorModel) {
//	funcName := "doViewUserParameter"
//	userParameterModel := repository.UserParameterModel{
//		UserID: sql.NullInt64{Int64: inputStruct.UserID},
//	}
//
//	userParameterModel, err = dao.UserParameterDAO.ViewUserParameter(serverconfig.ServerAttribute.DBConnection, userParameterModel)
//	if err.Error != nil {
//		return
//	}
//
//	if userParameterModel.ID.Int64 == 0 {
//		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
//		return
//	}
//
//	result = reformatDAOtoDTO(userParameterModel)
//	return
//}

//func reformatDAOtoDTO(categoryModel repository.UserParameterModel) out.ViewDetailUserParameterDTOOut {
//	var parameterValue map[string]string
//	_ = json.Unmarshal([]byte(categoryModel.ParameterValue.String), &parameterValue)
//
//	temp := out.ViewDetailUserParameterDTOOut{
//		ID:             categoryModel.ID.Int64,
//		UserID:         categoryModel.UserID.Int64,
//		ParameterValue: parameterValue,
//		CreatedBy:      categoryModel.CreatedBy.Int64,
//		UpdatedAt:      categoryModel.UpdatedAt.Time,
//	}
//
//	return temp
//}

//func (input userParameterService) validateView(inputStruct *in.UserParameterDTOIn) errorModel.ErrorModel {
//	return inputStruct.ValidateViewUserParameter()
//}
