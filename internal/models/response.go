package models

type Pagination struct {
	Page      int  `json:"page"`
	Limit     int  `json:"limit"`
	PrevPage  *int `json:"prev_page"`
	NextPage  *int `json:"next_page"`
	TotalPage int  `json:"total_page"`
}

type StandardResponse struct {
	Message    string      `json:"message"`
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}
