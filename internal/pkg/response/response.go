package response

var (
	MissingTeamIdCode       = "MISSING_TEAM_ID"
	ForbiddenCode           = "FORBIDDEN"
	IncorrectInputCode      = "INCORRECT_INPUT"
	InvalidInputCode        = "INVALID_INPUT"
	InternalServerErrorCode = "INTERNAL_SERVER_ERROR"
	NotAuthorizedCode       = "NOT_AUTHORIZED"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}
