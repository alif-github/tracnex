package CustomerContactService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

func (input customerContactService) insertCustomerContactForUpdateCustomer(tx *sql.Tx, inputStruct []in.CustomerContactRequest, contextModel *applicationModel.ContextModel, timeNow time.Time, isForUpdate bool) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var detailResponses out.CustomerErrorResponse

	for i := 0; i < len(inputStruct); i++ {
		err, detailResponses = input.ValidateInsertBulkCustomerContact(&inputStruct[i], contextModel)
		if err.Error != nil {
			output = detailResponses
			return
		}
	}

	output, dataAudit, err = input.InsertBulkCustomerContactFromCustomer(tx, inputStruct, contextModel, timeNow, true, false)
	return
}
