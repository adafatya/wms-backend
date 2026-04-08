package models

type StandardResponse struct {
	Message    string      `json:"message"`
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}
