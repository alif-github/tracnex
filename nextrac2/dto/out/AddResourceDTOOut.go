package out

// UserAuthDetail merupakan hasil response dari auth untuk api get user from internal.
type UserAuthDetail struct {
	AliasName         string                `json:"alias_name"`
	UserId            int64                 `json:"user_id"`
	Username          string                `json:"username"`
	FirstName         string                `json:"first_name"`
	LastName          string                `json:"last_name"`
	Email             string                `json:"email"`
	Phone             string                `json:"phone"`
	ClientId          string                `json:"client_id"`
	MaxAuthFail       int64                 `json:"max_auth_fail"`
	Locale            string                `json:"locale"`
	ResourceId        string                `json:"resource_id"`
	UserStatus        int64                 `json:"user_status"`
	SignatureKey      string                `json:"signature_key"`
	GrantTypes        string                `json:"grant_types"`
	Scope             string                `json:"scope"`
	IPWhitelist       string                `json:"ip_whitelist"`
	RedirectUri       string                `json:"redirect_uri"`
	UserInformation   AdditionalInformation `json:"user_information"`
	ClientInformation AdditionalInformation `json:"client_information"`
}

type AdditionalInformation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type APIResponseAddResource struct {
	Nexsoft NexsoftMessageAddResource `json:"nexsoft"`
}

type NexsoftMessageAddResource struct {
	Header  HeaderAddResource  `json:"header"`
	Payload PayloadAddResource `json:"payload"`
}

type PayloadAddResource struct {
	Status StatusResponseAddResource `json:"status"`
	Data   PayloadDataAddResource    `json:"data"`
	Other  interface{}               `json:"other"`
}

type StatusResponseAddResource struct {
	Success bool     `json:"success"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Detail  []string `json:"detail"`
}

type PayloadDataAddResource struct {
	Meta    interface{}    `json:"meta"`
	Content UserAuthDetail `json:"content"`
}

type HeaderAddResource struct {
	RequestID string `json:"request_id"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// Check Resource
type APIResponseCheckResource struct {
	Nexsoft NexsoftMessageCheckResource `json:"nexsoft"`
}

type NexsoftMessageCheckResource struct {
	Header  HeaderCheckResource  `json:"header"`
	Payload PayloadCheckResource `json:"payload"`
}

type PayloadCheckResource struct {
	Status StatusResponseCheckResource `json:"status"`
	Data   PayloadDataCheckResource    `json:"data"`
	Other  interface{}                 `json:"other"`
}

type StatusResponseCheckResource struct {
	Success bool     `json:"success"`
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Detail  []string `json:"detail"`
}

type PayloadDataCheckResource struct {
	Meta    interface{}          `json:"meta"`
	Content contentCheckResource `json:"content"`
}

type HeaderCheckResource struct {
	RequestID string `json:"request_id"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type contentCheckResource struct {
	IsExist               bool           `json:"is_exist"`
	AdditionalInformation UserAuthDetail `json:"additional_information"`
}
