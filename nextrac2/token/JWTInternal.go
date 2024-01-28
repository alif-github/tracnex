package token

import "github.com/dgrijalva/jwt-go"

type PayloadJWTInternal struct {
	Locale     string `json:"locale"`
	ClientID   string `json:"cid"`
	Resource   string `json:"resource"`
	Version    string `json:"version"`
	UserClient string `json:"user_client"`
	jwt.StandardClaims
}
