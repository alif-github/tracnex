package out

type GetSessionDTOOut struct {
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	Username       string              `json:"username"`
	Role           string              `json:"role"`
	Locale         string              `json:"locale"`
	UserID         int64               `json:"user_id"`
	IdCard         string              `json:"id_card"`
	Position       string              `json:"position"`
	Department     string              `json:"department"`
	IsHaveMember   bool                `json:"is_have_member"`
	Currency       string              `json:"currency"`
	PlatformDevice string              `json:"platform_device"`
	Scope          string              `json:"scope"`
	CurrentTime    string              `json:"current_time"`
	Permission     map[string][]string `json:"permission"`
	Menu           []ParentMenu        `json:"menu_parent"`
}

type GetSessionDateTimeDTOOut struct {
	CurrentTime string `json:"current_time"`
}

type ParentMenu struct {
	ID              int64         `json:"id"`
	Name            string        `json:"name"`
	EnName          string        `json:"en_name"`
	Sequence        int64         `json:"sequence"`
	IconName        string        `json:"icon_name"`
	Background      string        `json:"background"`
	AvailableAction string        `json:"available_action"`
	MenuCode        string        `json:"menu_code"`
	ServiceMenu     []ServiceMenu `json:"menu_service"`
}

type ServiceMenu struct {
	ID              int64      `json:"id"`
	ParentMenuID    int64      `json:"parent_menu_id"`
	Name            string     `json:"name"`
	EnName          string     `json:"en_name"`
	Sequence        int64      `json:"sequence"`
	IconName        string     `json:"icon_name"`
	Background      string     `json:"background"`
	MenuCode        string     `json:"menu_code"`
	AvailableAction string     `json:"available_action"`
	MenuItem        []MenuItem `json:"menu_item"`
}

type MenuItem struct {
	ID              int64  `json:"id"`
	ServiceMenuID   int64  `json:"service_menu_id"`
	Name            string `json:"name"`
	EnName          string `json:"en_name"`
	Sequence        int64  `json:"sequence"`
	IconName        string `json:"icon_name"`
	Background      string `json:"background"`
	URL             string `json:"url"`
	MenuCode        string `json:"menu_code"`
	AvailableAction string `json:"available_action"`
}

type MenuList struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	EnName          string `json:"en_name"`
	MenuCode        string `json:"menu_code"`
	AvailableAction string `json:"available_action"`
}

type ParentMenuList struct {
	ID              int64             `json:"id"`
	Name            string            `json:"name"`
	EnName          string            `json:"en_name"`
	MenuCode        string            `json:"menu_code"`
	IconName        string            `json:"icon_name"`
	AvailableAction string            `json:"available_action"`
	MenuService     []ServiceMenuList `json:"items"`
}

type ServiceMenuList struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	EnName          string         `json:"en_name"`
	MenuCode        string         `json:"menu_code"`
	IconName        string         `json:"icon_name"`
	AvailableAction string         `json:"available_action"`
	MenuItem        []ItemMenuList `json:"items"`
}

type ItemMenuList struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	EnName          string `json:"en_name"`
	MenuCode        string `json:"menu_code"`
	IconName        string `json:"icon_name"`
	AvailableAction string `json:"available_action"`
}

type NewMenuItemList struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	EnName          string `json:"en_name"`
	MenuCode        string `json:"menu_code"`
	AvailableAction string `json:"available_action"`
	UpdatedAt       string `json:"updated_at"`
}

type ViewDashboardResponse struct {
	TotalCustomerActive int64                `json:"total_customer_active"`
	TotalLicense        ViewLicenseDashboard `json:"total_license"`
}

type ViewLicenseDashboard struct {
	Total  int64                        `json:"total"`
	Detail []DetailViewLicenseDashboard `json:"detail"`
}

type DetailViewLicenseDashboard struct {
	Month int64 `json:"month"`
	Year  int64 `json:"year"`
	Total int64 `json:"total"`
}
