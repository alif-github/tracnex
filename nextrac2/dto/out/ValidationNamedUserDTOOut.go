package out

import "time"

type ValidationNamedUserResponse struct {
	Status           string    `json:"status"`
	ProductValidFrom time.Time `json:"product_valid_from"`
	ProductValidThru time.Time `json:"product_valid_thru"`
}
