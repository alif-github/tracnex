package constanta

var ListDepartmentNexsoft = []string{
	DepartmentDeveloper,
	DepartmentQAQC,
}

var ListTrackerQA = []string{
	TrackerAuto,
	TrackerManual,
}

var StatusAllowedDeveloper = []string{
	StatusNew, StatusReadyToDev, StatusInProgress, StatusCompleteDev,
	StatusNeedMoreReq, StatusReOpen, StatusClosed,
}

var StatusAllowedQA = []string{
	StatusNew, StatusReadyToTest, StatusInTest, StatusDoneTesting, StatusReOpen, StatusClosed, StatusReleased,
}
