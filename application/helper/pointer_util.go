package helper

func GetStringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StringPtr(s string) *string {
	return &s
}
