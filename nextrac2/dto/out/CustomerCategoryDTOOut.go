package out

import "time"

type CustomerCategoryResponse struct {
	ID                   int64     `json:"id"`
	CustomerCategoryID   string    `json:"customer_category_id"`
	CustomerCategoryName string    `json:"customer_category_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            int64     `json:"updated_by"`
	UpdatedName          string    `json:"updated_name"`
}

type CustomerCategoryDetailResponse struct {
	ID                   int64     `json:"id"`
	CustomerCategoryID   string    `json:"customer_category_id"`
	CustomerCategoryName string    `json:"customer_category_name"`
	CreatedBy            int64     `json:"created_by"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedBy            int64     `json:"updated_by"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedName          string    `json:"updated_name"`
}
