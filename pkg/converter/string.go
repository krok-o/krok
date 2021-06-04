package converter

// ToPointer takes a value and creates a pointer out from it.
func ToPointer(v string) *string {
	return &v
}
