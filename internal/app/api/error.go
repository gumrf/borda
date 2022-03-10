package api

type ErrorResponse struct {
	Errors []*APIError
}

type APIError struct {
	Status string
	Code   string
	Title  string
	Detail string
	Source ErrorSource `json:"omitempty"`
}

type ErrorSource struct {
	Parameter Parameter
}

type Parameter string
