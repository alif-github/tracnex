package out

type RegisterNamedUserResponse struct {
	UserRegistrationID int64 `json:"user_registration_id"`
}

type CheckEmailAndPhoneBeforeInsertResponse struct {
	Status  Status                       `json:"status"`
	Content []CheckEmailAndPhoneResponse `json:"content"`
}

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CheckEmailAndPhoneResponse struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	AuthUserId int64  `json:"auth_user_id"`
	Message    string `json:"message"`
}
