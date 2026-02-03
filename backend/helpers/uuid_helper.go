package helpers

import (
	"fmt"

	"github.com/google/uuid"
)

func ParseUUIDPointer(id string) *uuid.UUID {
	u := uuid.MustParse(id)
	return &u
}

// ParsePtrToUUIDPtr parses a string pointer into a *uuid.UUID.
// Returns nil if the input is nil or empty.
// Returns an error if the input is not a valid UUID.
func ParsePtrToUUIDPtr(input *string) (*uuid.UUID, error) {
	if input == nil || *input == "" {
		return nil, nil // Treat nil or empty strings as nil UUID
	}

	parsed, err := uuid.Parse(*input)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: '%s'", *input)
	}
	return &parsed, nil
}

// ParseUUIDToPtr parsing valid uuid string to uuid.UUID and return pointer to it
func ParseUUIDToPtr(input string) (*uuid.UUID, error) {
	u, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// MustParseUUIDToPtr parse a valid uuid string to uuid.UUID and return pointer to it
func MustParseUUIDToPtr(input string) *uuid.UUID {
	u := uuid.MustParse(input)
	return &u
}

// UUIDToPtr converts a UUID value to a pointer
func UUIDToPtr(value uuid.UUID) *uuid.UUID {
	return &value
}

// PtrToUUID dereferences a UUID pointer, returning a default value if nil
func PtrToUUID(ptr *uuid.UUID, defaultValue uuid.UUID) uuid.UUID {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// IsPtrEqualsToUUID compares a UUID pointer with a UUID value
func IsPtrEqualsToUUID(ptr *uuid.UUID, value uuid.UUID) bool {
	if ptr == nil {
		return false
	}
	return *ptr == value
}

// CreateUUIDPtr creates a new UUID and returns it as a pointer
func CreateUUIDPtr() *uuid.UUID {
	id := uuid.New()
	return &id
}
