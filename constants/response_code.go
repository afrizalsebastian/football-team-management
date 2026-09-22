package constants

import "net/http"

type ResponseCode struct {
	message  string
	code     int
	httpCode int
}

func (r *ResponseCode) GetMessage() string {
	return r.message
}

func (r *ResponseCode) GetCode() int {
	return r.code
}

func (r *ResponseCode) GetHttpCode() int {
	return r.httpCode
}

func registerResponseCode(message string, code int, httpCode int) *ResponseCode {
	return &ResponseCode{
		message:  message,
		code:     code,
		httpCode: httpCode,
	}
}

var (

	// example
	SuccessHello = registerResponseCode("Success -- Hello from App", 200999, http.StatusOK)

	// 2XX
	SuccessDefault        = registerResponseCode("Success", 200000, http.StatusOK)
	SuccesssWithEmptyList = registerResponseCode("Success with empty list data", 200001, http.StatusOK)
	CreatedDefault        = registerResponseCode("Success", 201000, http.StatusCreated)
	AcceptDefault         = registerResponseCode("Accpeted", 202000, http.StatusAccepted)

	// 4XX
	BadRequestDefault = registerResponseCode("Bad request", 400000, http.StatusBadRequest)

	// 5XX
	InternalServerError = registerResponseCode("Internal Server Error", 500000, http.StatusInternalServerError)
)
