package out

import "time"

type WhiteListDeviceResponse struct {
	ID          int64     `json:"id"`
	Device      string    `json:"device"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedName string    `json:"updated_name"`
}