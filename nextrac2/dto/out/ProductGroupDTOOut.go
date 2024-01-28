package out

import "time"

type ProductGroupResponse struct {
	ID               int64     `json:"id"`
	ProductGroupName string    `json:"product_group_name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        int64     `json:"updated_by"`
	UpdatedName      string    `json:"updated_name"`
}

type ProductGroupDetailResponse struct {
	ID               int64     `json:"id"`
	ProductGroupName string    `json:"product_group_name"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        int64     `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        int64     `json:"updated_by"`
	UpdatedName      string    `json:"updated_name"`
}