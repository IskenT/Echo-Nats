package api

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

const ItemNotFoundMessage = "errors.item.notFound"
const ItemNotFoundCode = 3

func NewErrorResponse(code int, message string, details ...interface{}) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}
