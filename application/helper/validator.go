package helper

import (
	"context"
	"regexp"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation if needed
	validate.RegisterValidation("alphanumeric", validateAlphanumeric)
}

func validateAlphanumeric(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return re.MatchString(field)
}

func getValidationErrorMessage(ctx context.Context, fe validator.FieldError) (int, string) {
	var (
		errCode          int
		validationErrMsg string
	)

	switch fe.Tag() {
	case "required":
		errCode = constants.RequiredValidationErr.GetCode()
		validationErrMsg = constants.RequiredValidationErr.GetMessageWithParam(fe.Field())
	case "gt":
		errCode = constants.GreaterThanValidationErr.GetCode()
		validationErrMsg = constants.GreaterThanValidationErr.GetMessageWithParam(fe.Field(), fe.Param())
	case "max":
		errCode = constants.MaxLengthValidationErr.GetCode()
		validationErrMsg = constants.MaxLengthValidationErr.GetMessageWithParam(fe.Field(), fe.Param())
	case "oneof":
		errCode = constants.OneOfValidationErr.GetCode()
		validationErrMsg = constants.OneOfValidationErr.GetMessageWithParam(fe.Field(), fe.Param())
	case "min":
		errCode = constants.MinLengthValidationErr.GetCode()
		validationErrMsg = constants.MinLengthValidationErr.GetMessageWithParam(fe.Field(), fe.Param())
	case "email":
		errCode = constants.EmailValidationErr.GetCode()
		validationErrMsg = constants.EmailValidationErr.GetMessageWithParam(fe.Field())
	default:
		errCode = constants.DefaultValidationErr.GetCode()
		validationErrMsg = constants.DefaultValidationErr.GetMessageWithParam(fe.Field())
	}

	return errCode, validationErrMsg
}

func ValidateParams(ctx context.Context, req interface{}) []api.ErrorsDetail {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	var fieldErrs []api.ErrorsDetail
	for _, fe := range err.(validator.ValidationErrors) {
		code, msg := getValidationErrorMessage(ctx, fe)
		fieldErrs = append(fieldErrs, api.ErrorsDetail{
			Field:   fe.Field(),
			Message: msg,
			Code:    code,
		})
	}

	return fieldErrs
}
