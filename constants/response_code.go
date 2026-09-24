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
	BadRequestDefault   = registerResponseCode("Bad request", 400000, http.StatusBadRequest)
	InvalidJerseyNumber = registerResponseCode("Jersery number already used by active player", 400001, http.StatusBadRequest)
	MatchDateInPast     = registerResponseCode("Match date in past. Use other date", 400002, http.StatusBadRequest)

	NotFoundDefault = registerResponseCode("Not Found", 404000, http.StatusNotFound)
	NotFoundTeam    = registerResponseCode("Not Found Team", 404001, http.StatusNotFound)

	// 5XX
	InternalServerError = registerResponseCode("Internal Server Error", 500000, http.StatusInternalServerError)
)
