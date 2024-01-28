package out

import "time"

//--------------Get List Response

type GetListPersonTitle struct {
	Nexsoft NexsoftListPersonTitleResponse `json:"nexsoft"`
}

type NexsoftListPersonTitleResponse struct {
	Header  Header                 `json:"header"`
	Payload PayloadListPersonTitle `json:"payload"`
}

type PayloadListPersonTitle struct {
	Status StatusResponse  `json:"status"`
	Data   ListPersonTitle `json:"data"`
	Other  interface{}     `json:"other"`
}

type ListPersonTitle struct {
	Meta    interface{}           `json:"meta"`
	Content []PersonTitleResponse `json:"content"`
}

//--------------View response

type ViewPersonTitle struct {
	Nexsoft NexsoftViewPersonTitleResponse `json:"nexsoft"`
}

type NexsoftViewPersonTitleResponse struct {
	Header  Header                 `json:"header"`
	Payload PayloadViewPersonTitle `json:"payload"`
}

type PayloadViewPersonTitle struct {
	Status StatusResponse      `json:"status"`
	Data   DataViewPersonTitle `json:"data"`
	Other  interface{}         `json:"other"`
}

type DataViewPersonTitle struct {
	Meta    interface{}         `json:"meta"`
	Content PersonTitleResponse `json:"content"`
}

//--------------Content

type PersonTitleResponse struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	EnDescription string    `json:"en_description"`
	Status        string    `json:"status"`
	CreatedBy     int64     `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
}
