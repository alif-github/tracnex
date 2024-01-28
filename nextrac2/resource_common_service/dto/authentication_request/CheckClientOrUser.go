package authentication_request

type CheckClientOrUser struct {
	Username	string		`json:"username"`
	Email		string		`json:"email"`
	Phone		string		`json:"phone"`
	ClientID	string		`json:"client_id"`
}
