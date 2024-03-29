package ReportService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input reportService) SetInitPaidSprint(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		tx   *sql.Tx
		errs error
	)

	defer func() {
		if errs != nil || err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	tx, errs = serverconfig.ServerAttribute.RedmineDBConnection.Begin()
	if errs != nil {
		return
	}

	rawSprint := "---\n- 20230806-20230820\n- 20230721-20230805\n- 20230706-20230720\n- 20230621-20230705\n- 20230606-20230620\n- 20230521-20230605\n- 20230506-20230520\n- 20230421-20230505\n- 20230406-20230420\n- 20230321-20230405\n- 20230306-20230320\n- 20230221-20230305\n- 20230206-20230220\n- 20230121-20230205\n- 20230106-20230120\n- 20221221-20230105\n- 20221206-20221220\n- 20221121-20221205\n- 20221106-20221120\n- 20221021-20221105\n- 20221006-20221020\n- 20220921-20221005\n- 20220906-20220920\n- 20220821-20220905\n- 20220806-20220820\n- 20220721-20220805\n- 20220706-20220720\n- 20220621-20220705\n- 20220606-20220620\n- 20220521-20220605\n- 20220506-20220520\n- 20220421-20220505\n- 20220406-20220420\n- 20220321-20220405\n- 20220306-20220320\n- 20220221-20220305\n- 20220206-20220220\n- 20220121-20220205\n- 20220106-20220120\n- 20211221-20220105\n- 20211206-20211220\n- 20211121-20211205\n- 20211106-20211120\n- 20211021-20211105\n- 20211006-20211020\n- 20210921-20211005\n- 20210906-20210920\n- 20210821-20210905\n- 20210806-20210820\n- 20210721-20210805\n- 20210706-20210720\n- 20210621-20210705\n- 20210606-20210620\n- 20210521-20210605\n- 20210506-20210520\n- 20210421-20210505\n- 20210406-20210420\n- 20210321-20210405\n- 20210306-20210320\n- 20210221-20210305\n- 20210206-20210220\n- 20210121-20210205\n- 20210106-20210120\n- 20201221-20210105\n- 20201206-20201220\n- 20201121-20201205\n- 20201106-20201120\n- 20201021-20201105\n- 20201005-20201020\n- 20200921-20201005\n- 20200906-20200920\n- 20200821-20200905\n- 20200806-20200820\n- 20200721-20200805\n- 20200706-20200720\n- 20200621-20200705\n- 20200606-20200620\n- 20200521-20200605\n- 20200506-20200520\n- 20200421-20200505\n- 20200406-20200420\n- 20200321-20200405\n- 20200306-20200320\n"
	rawPaid := "---\n- UNPAID\n- REJECTED\n- PAID-2020-03\n- PAID-2020-04\n- PAID-2020-05\n- PAID-2020-06\n- PAID-2020-07\n- PAID-2020-08\n- PAID-2020-09\n- PAID-2020-10\n- PAID-2020-11\n- PAID-2020-12\n- PAID-2021-01\n- PAID-2021-02\n- PAID-2021-03\n- PAID-2021-04\n- PAID-2021-05\n- PAID-2021-06\n- PAID-2021-07\n- PAID-2021-08\n- PAID-2021-09\n- PAID-2021-10\n- PAID-2021-11\n- PAID-2021-12\n- PAID-2022-01\n- PAID-2022-02\n- PAID-2022-03\n- PAID-2022-04\n- PAID-2022-05\n- PAID-2022-06\n- PAID-2022-07\n- PAID-2022-08\n- PAID-2022-09\n- PAID-2022-10\n- PAID-2022-11\n- PAID-2022-12\n- PAID-2023-01\n- PAID-2023-02\n- PAID-2023-03\n- PAID-2023-04\n- PAID-2023-05\n- PAID-2023-06\n- PAID-2023-07\n"

	//--- Update Sprint
	err = dao.RedmineDAO.UpdateCustomFieldsOnRedmine(tx, constanta.IDSprintOnRedmineCustomFields, rawSprint)
	if err.Error != nil {
		return
	}

	//--- Update Paid
	err = dao.RedmineDAO.UpdateCustomFieldsOnRedmine(tx, constanta.IDPaymentOnRedmineCustomFields, rawPaid)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}
