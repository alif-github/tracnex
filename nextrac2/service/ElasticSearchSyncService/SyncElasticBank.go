package ElasticSearchSyncService

//func (input elasticSearchSyncService) GetSyncElasticBankChildTask() backgroundJobModel.ChildTask {
//	return backgroundJobModel.ChildTask{
//		Group: constanta.JobProcessSynchronizeGroup,
//		Type:  constanta.JobProcessElasticType,
//		SocketID:  constanta.JobProcessSyncElasticBank,
//		Data: backgroundJobModel.BackgroundServiceModel{
//			SearchByParam: nil,
//			IsCheckStatus: false,
//			CreatedBy:     0,
//			Data:          nil,
//		},
//		GetCountData: dao.BankDAO.GetCountBank,
//		DoJob:        input.syncBank,
//	}
//}
//
//func (input elasticSearchSyncService) DoSyncElasticBank(contextModel applicationModel.ContextModel) {
//	timeNow := time.Now()
//	task := input.GetSyncElasticBankChildTask()
//	job := service.GetJobProcess(task, contextModel, timeNow)
//	input.ServiceWithBackgroundProcess(serverconfig.ServerAttribute.DBConnection, false, repository.JobProcessModel{}, job, task, contextModel)
//}
//
//func (input elasticSearchSyncService) syncBank(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel) (err errorModel.ErrorModel) {
//	input.clearIndex(serverconfig.ServerAttribute.ElasticClient, input.addPrefixIndexName(dao.BankDAO.ElasticSearchIndex))
//	err = input.syncElasticData(db, nil, childJob, dao.BankDAO.GetListBank, func(data interface{}) errorModel.ErrorModel {
//		bankModel := data.(repository.BankModel)
//		idStr := strconv.Itoa(int(bankModel.ID.Int64))
//		dao.BankDAO.DoInsertAtElasticSearch(idStr, repository.BankElasticModel{
//			ID:        bankModel.ID.Int64,
//			SocketID:      bankModel.SocketID.String,
//			Status:    bankModel.Status.String,
//			CreatedBy: bankModel.CreatedBy.Int64,
//		})
//		return errorModel.GenerateNonErrorModel()
//	})
//
//	return
//}
