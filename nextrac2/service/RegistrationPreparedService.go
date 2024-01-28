package service

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

type RegistrationPrepared struct {
	IdResourceAllowed []int64
}

func (input RegistrationPrepared) CheckIsClientTypeExist(inputClientTypeID int64, preparedError errorModel.ErrorModel) (err errorModel.ErrorModel) {

	result, err := dao.ClientTypeDAO.CheckClientTypeByID(serverconfig.ServerAttribute.DBConnection, &repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: inputClientTypeID},
	})
	if err.Error != nil {
		return
	}

	if result.ID.Int64 == 0 {
		err = preparedError
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}