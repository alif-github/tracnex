package UserParameterService

//func (input userParameterService) UpdateUserParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
//	var inputStruct in.ParameterRequest
//
//	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
//	if err.Error != nil {
//		return
//	}
//
//	_, err = input.ServiceWithDataAuditPreparedByService("UpdateUserParameter", inputStruct, contextModel, input.doUpdateUserParameter, func(_ interface{}, _ applicationModel.ContextModel) {})
//	if err.Error != nil {
//		return
//	}
//
//	output.Status = out.StatusResponse{
//		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
//		MessageAuth: GenerateI18NMessage("SUCCESS_UPDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
//	}
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input userParameterService) doUpdateUserParameter(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
//	funcName := "doUpdateUserParameter"
//
//	inputStruct := inputStructInterface.(in.ParameterRequest)
//	var userParameterDB repository.UserParameterModel
//	var userParameterID int64
//
//	userParameterModel := repository.UserParameterModel{
//		UserID:         sql.NullInt64{Int64: inputStruct.ID},
//		ParameterValue: sql.NullString{String: util.StructToJSON(inputStruct.Value)},
//		CreatedAt:      sql.NullTime{Time: timeNow},
//		CreatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
//		UpdatedAt:      sql.NullTime{Time: timeNow},
//		UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
//		UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
//	}
//
//	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
//	if isOnlyHaveOwnAccess {
//		userParameterModel.CreatedBy.Int64 = userID
//	}
//
//	userParameterDB, err = dao.UserParameterDAO.GetUserParameterForUpdate(tx, userParameterModel)
//	if err.Error != nil {
//		return
//	}
//
//	if userParameterDB.ID.Int64 == 0 {
//		userParameterModel.CreatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
//		userParameterID, err = dao.UserParameterDAO.InsertUserParameter(tx, userParameterModel)
//		if err.Error != nil {
//			err = input.CheckDuplicateError(err)
//			return
//		}
//		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.UserParameterDAO.TableName, userParameterID, contextModel.LimitedByCreatedBy)...)
//	} else {
//		if userParameterDB.UpdatedAt.Time != inputStruct.UpdatedAt {
//			err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.UserParameter)
//			return
//		}
//
//		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserParameterDAO.TableName, userParameterDB.ID.Int64, contextModel.LimitedByCreatedBy)...)
//		err = dao.UserParameterDAO.UpdateUserParameter(tx, userParameterModel, timeNow)
//		if err.Error != nil {
//			err = input.CheckDuplicateError(err)
//			return
//		}
//	}
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input userParameterService) validateUpdate(inputStruct *in.UserParameterDTOIn) errorModel.ErrorModel {
//	return inputStruct.ValidateUpdateUserParameter()
//}
