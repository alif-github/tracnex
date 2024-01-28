package token

import "github.com/dgrijalva/jwt-go"

type PayloadJWTToken struct {
	ClientID string `json:"cid"`
	Resource string `json:"resource"`
	Scope    string `json:"scope"`
	Locale   string `json:"locale"`
	jwt.StandardClaims
}
