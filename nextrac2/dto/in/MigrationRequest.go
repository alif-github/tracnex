package in

import "nexsoft.co.id/nextrac2/model/errorModel"

type ResetMigrationRequest struct {
	ID    []string `json:"id"`
	Reset bool
}

func (input *ResetMigrationRequest) ValidateReset() (err errorModel.ErrorModel) {
	if len(input.ID) < 1 {
		// reset all migrations
		input.Reset = true
	}

	return
}
