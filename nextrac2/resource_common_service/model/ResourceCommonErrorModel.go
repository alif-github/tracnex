package model

type ResourceCommonErrorModel struct {
	Code     int
	Error    error
	CausedBy error
}
