package cryptoModel

type EncryptLicenseResponseModel struct {
	MessageCode          string                `json:"MessageCode"`
	Message              string                `json:"Message"`
	Notification         string                `json:"Notification"`
	ProductSignature     string                `json:"ProductSignature"`
	ProductEncrypt       string                `json:"ProductEncrypt"`
	ProductKey           string                `json:"ProductKey"`
	ProductConfiguration EncryptLicenseRequest `json:"ProductConfiguration"`
}
