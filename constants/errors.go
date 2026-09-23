package constants

import "errors"

var (
	InvalidUUIDValue  = errors.New("Invalid UUID value")
	InvalidFieldValue = errors.New("Invalid Field Value")
	DuplicateRow      = errors.New("Duplicate Row in DB")
)
