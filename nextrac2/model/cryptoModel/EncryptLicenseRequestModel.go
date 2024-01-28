package cryptoModel

type EncryptLicenseRequest struct {
	SignatureKey     string `json:"signatureKey"`
	ProductSignature string `json:"productSignature"`
	ClientSecret     string `json:"clientSecret"`
	ClientId         string `json:"clientId"`
	EncryptKey       string `json:"encryptKey"`
	ProductKey       string `json:"productKey"`
	HardwareId       string `json:"hardwareId"`
	ProductId        string `json:"productId"`
}

type EncryptLicenseRequestModel struct {
	SignatureKey      string                         `json:"signature_key"`
	ClientSecret      string                         `json:"client_secret"`
	Hwid              string                         `json:"hwid"`
	LicenseConfigData JSONFileActivationLicenseModel `json:"license_config_data"`
}

type ReEncryptLicenseModel struct {
	EncryptLicenseRequestModel
	ProductKey string `json:"product_key"`
}
