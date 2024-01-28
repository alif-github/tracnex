// this is a common file that will be used for all the common 
// definitions and structures in the platform

package includes

type Authorisation struct {
	Authz string
}

type User_details struct {
	User_id   int64
	User_name string
}

type User struct {
	LoggedInUser User_details
	AuthzData    Authorisation
}

type Rc struct {
	Code    string
	Message string
}

type nexsoft_error struct {
	Err error // this is golang error type
	Rc  Rc    // nexsoft standard error returned from DB

}
