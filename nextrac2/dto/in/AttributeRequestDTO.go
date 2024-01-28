package in

//---------- Attribute request for registration client service
type ClientRequestAttributeRequest struct {
	ClientTypeID 	int64        					`json:"client_type_id"`
	ClientName   	string       					`json:"client_name"`
	SocketID     	string       					`json:"socket_id"`
	CompanyData  	[]CompanyDataAttributeRequest 	`json:"company_data"`
}

type CompanyDataAttributeRequest struct {
	CompanyID		string		 					`json:"company_id"`
	BranchData     	[]BranchDataAttributeRequest 	`json:"branch_data"`
}

type BranchDataAttributeRequest struct {
	BranchID		string			`json:"branch_id"`
}

//---------- Attribute request for registration pkce
type AttributeRequestRegistPKCE struct {
	ParentClientID 		string 	`json:"parent_client_id"`
	ClientTypeID   		int64  	`json:"client_type_id"`
	CompanyID      		string 	`json:"company_id"`
	BranchID       		string 	`json:"branch_id"`
	Username       		string 	`json:"username"`
	Password       		string 	`json:"password"`
	FirstName      		string 	`json:"first_name"`
	LastName       		string 	`json:"last_name"`
	Email          		string 	`json:"email"`
	Phone          		string 	`json:"phone"`
}

//---------- Attribute request for add resource
type AttributeRequestAddResource struct {
	ClientTypeID   		int64  	`json:"client_type_id"`
	ClientID			string	`json:"client_id"`
}

//---------- Attribute request for error registered client
type AttributeRequestErrorRegisteredClient struct {
	CompanyID		string	`json:"company_id"`
	BranchID		string	`json:"branch_id"`
}