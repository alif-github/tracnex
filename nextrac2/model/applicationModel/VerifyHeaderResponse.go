package applicationModel

type VerifyHeaderResponse struct {
	RedirectURI string
	Token       string
	Refresh     string
	State       string
	Code        string
}
