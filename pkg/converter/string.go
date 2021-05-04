package converter

// ToValue takes a pointer to a string value and returns the value.
func ToValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// ToPointer takes a value and creates a pointer out from it.
func ToPointer(v string) *string {
	return &v
}
