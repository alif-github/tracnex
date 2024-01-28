package out

import "time"

type ComponentResponse struct {
	ID                   int64     `json:"id"`
	ComponentName 		 string    `json:"component_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            int64     `json:"updated_by"`
	UpdatedName          string    `json:"updated_name"`
}