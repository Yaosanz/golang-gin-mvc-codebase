package helpers

// BoolToPtr creates a pointer to the given boolean value
//
//	Usage:
//	- use this helper if you need to create a pointer from a boolean value:
//	Params:
//	- value: bool
//	Returns:
//	- *bool
func BoolToPtr(value bool) *bool {
	return &value
}

// PtrToBool dereferences a boolean pointer, returning the default value if nil
//
//	Usage:
//	- use this helper to dereference a Pointer to Boolean
//	Params:
//	- ptr: *bool
//	- defaultValue: bool
//	Returns:
//	- bool
func PtrToBool(ptr *bool, defaultValue bool) bool {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// IsPtrEqualsToBool compares a boolean pointer (*bool) with a boolean value
//
//	Usage:
//	- use this helper If you need to compare a pointer boolean value (*bool) with a regular boolean value (bool):
//	Params:
//	- ptr: *bool
//	- value: bool
//	Returns:
//	- bool
func IsPtrEqualsToBool(ptr *bool, value bool) bool {
	if ptr == nil {
		return false // Nil pointer dianggap tidak sama
	}
	return *ptr == value
}
