package internal

import "net/http"

type ResponseCode struct {
	Message  string
	Code     int
	HttpCode int
}

func NewResponseCode(message string, code int, httpCode int) ResponseCode {
	return ResponseCode{
		Message:  message,
		Code:     code,
		HttpCode: httpCode,
	}
}

var (
	SuccessHello = NewResponseCode("Success -- Hello from App", 200999, http.StatusOK)

	//
	SuccessDefault = NewResponseCode("Success", 200000, http.StatusOK)
	CreatedDefault = NewResponseCode("Success", 201000, http.StatusCreated)
	AcceptDefault  = NewResponseCode("Accpeted", 202000, http.StatusAccepted)

	BadRequestDefault = NewResponseCode("Bad request", 400000, http.StatusBadRequest)

	InternalServerError = NewResponseCode("Internal Server Error", 500000, http.StatusInternalServerError)
)
