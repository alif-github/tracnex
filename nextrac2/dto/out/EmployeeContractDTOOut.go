package out

import "time"

type GetListEmployeeContractResponse struct {
	ID          int64     `json:"id"`
	ContractNo  string    `json:"contract_no"`
	Information string    `json:"information"`
	EmployeeID  int64     `json:"employee_id"`
	FromDate    time.Time `json:"from_date"`
	ThruDate    time.Time `json:"thru_date"`
	CreatedName string    `json:"created_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedName string    `json:"updated_name"`
	UpdatedAt   time.Time `json:"updated_at"`
}
