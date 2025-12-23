package helpers

// Integer constraint mencakup semua tipe integer di Go
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// =============== helper to convert any integer to int =============== //

// ToInt convert any integer to int
func ToInt[T Integer](v T) int {
	return int(v)
}

// ToIntPtr convert any integer to *int
func ToIntPtr[T Integer](v T) *int {
	val := int(v)
	return &val
}

// IntToPtr convert any int ke *int
func IntToPtr(v int) *int {
	return &v
}

// =============== helper to convert any integer to int8 =============== //

// ToInt8 convert any integer to int8
func ToInt8[T Integer](v T) int8 {
	return int8(v)
}

// ToInt8Ptr convert any integer to *int8
func ToInt8Ptr[T Integer](v T) *int8 {
	val := int8(v)
	return &val
}

// Int8ToPtr convert any int8 ke *int8
func Int8ToPtr(v int8) *int8 {
	return &v
}

// =============== helper to convert any integer to int16 =============== //

// ToInt16 convert any integer to int16
func ToInt16[T Integer](v T) int16 {
	return int16(v)
}

// ToInt16Ptr convert any integer to *int16
func ToInt16Ptr[T Integer](v T) *int16 {
	val := int16(v)
	return &val
}

// Int16ToPtr convert any int16 ke *int16
func Int16ToPtr(v int16) *int16 {
	return &v
}

// =============== helper to convert any integer to int32 =============== //

// ToInt32 convert any integer to int32
func ToInt32[T Integer](v T) int32 {
	return int32(v)
}

// ToInt32Ptr convert any integer to *int32
func ToInt32Ptr[T Integer](v T) *int32 {
	val := int32(v)
	return &val
}

// Int32ToPtr convert any int32 ke *int32
func Int32ToPtr(v int32) *int32 {
	return &v
}

// =============== helper to convert any integer to int64 =============== //

// ToInt64 convert any integer to int64
func ToInt64[T Integer](v T) int64 {
	return int64(v)
}

// ToInt64Ptr convert any integer to *int64
func ToInt64Ptr[T Integer](v T) *int64 {
	val := int64(v)
	return &val
}

// Int64ToPtr convert any int64 ke *int64
func Int64ToPtr(v int64) *int64 {
	return &v
}
