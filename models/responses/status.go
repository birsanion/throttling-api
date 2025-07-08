package responses

type SuccessRespose struct {
	Success bool `json:"success"`
}

func NewSuccessResponse(success bool) SuccessRespose {
	return SuccessRespose{Success: success}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Error: msg}
}
