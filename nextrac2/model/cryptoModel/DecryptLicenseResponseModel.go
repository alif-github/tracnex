package cryptoModel

type DecryptLicenseResponseModel struct {
	MessageCode      string                         `json:"MessageCode"`
	Message          string                         `json:"Message"`
	Notification     string                         `json:"Notification"`
	ProductSignature string                         `json:"ProductSignature"`
	ProductEncrypt   string                         `json:"ProductEncrypt"`
	ProductKey       string                         `json:"ProductKey"`
	Configuration    JSONFileActivationLicenseModel `json:"Configuration"`
}
