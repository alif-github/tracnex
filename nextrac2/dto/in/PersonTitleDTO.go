package in

type PersonTitleRequest struct {
	AbstractDTO
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	EnDescription string `json:"en_description"`
}