package helper

func GetStringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GetFloat64PtrValue(v *float64) float64 {
	if v == nil {
		return 0
	}

	return *v
}

func GetIntPtrValue(v *int) int {
	if v == nil {
		return 0
	}

	return *v
}

func StringPtr(s string) *string {
	return &s
}

func Float64Ptr(s float64) *float64 {
	return &s
}

func IntPtr(s int) *int {
	return &s
}
