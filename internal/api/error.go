package api

type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
}

type ErrorObject struct {
	Status string        `json:"status"`
	Code   string        `json:"code"`
	Title  string        `json:"title,omitempty"`
	Detail string        `json:"detail,omitempty"`
	Source *SourceObject `json:"source,omitempty"`
}

type SourceObject struct {
	Parameter Parameter `json:"parametr"`
}

type Parameter string

func NewAPIErrorResponse(errors ...ErrorObject) *ErrorResponse {
	return &ErrorResponse{
		Errors: errors,
	}
}
