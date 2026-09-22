package constants

import "fmt"

type ValidationErr struct {
	code int
	msg  string
}

func (v *ValidationErr) GetMessageWithParam(params ...interface{}) string {
	return fmt.Sprintf(v.msg, params...)
}

func (v *ValidationErr) GetCode() int {
	return v.code
}

func regValidationErr(code int, msg string) *ValidationErr {
	return &ValidationErr{
		code: code,
		msg:  msg,
	}
}

var (
	DefaultValidationErr     = regValidationErr(400100, "This field '%s' has an invalid value")
	RequiredValidationErr    = regValidationErr(400101, "This field '%s' is required")
	GreaterThanValidationErr = regValidationErr(400102, "This field '%s' must be greater than %s")
	MaxLengthValidationErr   = regValidationErr(400103, "This field '%s' must be less than or equal to %s characters")
	MinLengthValidationErr   = regValidationErr(400104, "This field '%s' must be at least %s characters")
	OneOfValidationErr       = regValidationErr(400105, "This field '%s' must be one of '%s'")
	EmailValidationErr       = regValidationErr(400107, "This field '%s' must be a valid email")
)
