package cryptoModel

type DecryptLicenseRequestModel struct {
	SignatureKey     string `json:"signatureKey"`
	ProductSignature string `json:"productSignature"`
	ClientId         string `json:"clientId"`
	ClientSecret     string `json:"clientSecret"`
	EncryptKey       string `json:"encryptKey"`
	HardwareId       string `json:"hardwareId"`
	ProductKey       string `json:"productKey"`
	ProductId        string `json:"productId"`
}
