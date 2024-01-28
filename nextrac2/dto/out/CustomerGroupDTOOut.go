package out

import "time"

type CustomerGroupResponse struct {
	ID                int64     `json:"id"`
	CustomerGroupID   string    `json:"customer_group_id"`
	CustomerGroupName string    `json:"customer_group_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedBy         int64     `json:"updated_by"`
	UpdatedName       string    `json:"updated_name"`
}

type CustomerGroupDetailResponse struct {
	ID                int64     `json:"id"`
	CustomerGroupID   string    `json:"customer_group_id"`
	CustomerGroupName string    `json:"customer_group_name"`
	CreatedBy         int64     `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedBy         int64     `json:"updated_by"`
	UpdatedAt         time.Time `json:"updated_at"`
	UpdatedName       string    `json:"updated_name"`
}
