package helper

import (
	"context"
	"regexp"
	"time"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation if needed
	validate.RegisterValidation("ddmmyyyy", validateDateFormat)
	validate.RegisterValidation("hhmm", validateTimeFormat)
	validate.RegisterValidation("goal", validateGoalMinute)
}

func validateDateFormat(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}
	_, err := time.Parse("02-01-2006", field)
	if err != nil {
		return false
	}
	return true
}

func validateTimeFormat(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}
	matched, _ := regexp.MatchString(`^([01]\d|2[0-3]):([0-5]\d)$`, field)
	return matched
}

func validateGoalMinute(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	matched, _ := regexp.MatchString(`^\d{1,3}:[0-5]\d(\+\d{1,3}:[0-5]\d)?$`, field)
	return matched
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
	case "ddmmyyyy":
		errCode = constants.DateValidationErr.GetCode()
		validationErrMsg = constants.DateValidationErr.GetMessageWithParam(fe.Field())
	case "hhmm":
		errCode = constants.TimeValidationErr.GetCode()
		validationErrMsg = constants.TimeValidationErr.GetMessageWithParam(fe.Field())
	case "nefield":
		errCode = constants.NotEqualFieldValidationErr.GetCode()
		validationErrMsg = constants.NotEqualFieldValidationErr.GetMessageWithParam(fe.Field(), fe.Param())
	case "goal":
		errCode = constants.GoalMinuteValidationErr.GetCode()
		validationErrMsg = constants.GoalMinuteValidationErr.GetMessageWithParam(fe.Field())
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
