package out

type RegisterOrRenewLicenseUserResponse struct {
	UserRegistrationDetailId int64  `json:"user_registration_detail_id"`
	ClientId                 string `json:"client_id"`
	AuthUserId               int64  `json:"auth_user_id"`
}
