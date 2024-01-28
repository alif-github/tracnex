package repository

type EmployeeNotification struct {
	EmployeeId     int64
	MemberIdList   []string
	FilterByIsRead bool
	IsRead         bool
}