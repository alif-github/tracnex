package out

import "time"

type GetListUserDTOOut struct {
	ID             int64     `json:"id"`
	ClientID       string    `json:"client_id"`
	AuthUserID     int64     `json:"auth_user_id"`
	Username       string    `json:"username"`
	Firstname      string    `json:"firstname"`
	Lastname       string    `json:"lastname"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	RoleID         string    `json:"role_id"`
	GroupID        string    `json:"group_id"`
	Locale         string    `json:"locale"`
	Status         string    `json:"status"`
	CreatedName    string    `json:"created_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PlatformDevice string    `json:"platform_device"`
}

type ViewUserDTOOut struct {
	Username        string    `json:"username"`
	Firstname       string    `json:"firstname"`
	Lastname        string    `json:"lastname"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Role            string    `json:"role"`
	GroupID         string    `json:"group_id"`
	IsAdmin         bool      `json:"is_admin"`
	Status          string    `json:"status"`
	StatusDefine    string    `json:"status_define"`
	IsDisableStatus bool      `json:"is_disable_status"`
	IsVerifyPhone   bool      `json:"is_verify_phone"`
	IsVerifyEmail   bool      `json:"is_verify_email"`
	CreatedBy       int64     `json:"created_by"`
	CreatedName     string    `json:"created_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedBy       int64     `json:"updated_by"`
	UpdatedName     string    `json:"updated_name"`
	UpdatedAt       time.Time `json:"updated_at"`
	PlatformDevice  string    `json:"platform_device"`
}
