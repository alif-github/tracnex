package out

import "time"

type ModuleResponse struct {
	ID                   int64     `json:"id"`
	ModuleName 			 string    `json:"module_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            int64     `json:"updated_by"`
	UpdatedName          string    `json:"updated_name"`
}
