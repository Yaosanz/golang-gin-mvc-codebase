package database

import "fmt"

// Error Custom error types for better error handling
type Error struct {
	Component string
	Type      string
	Err       error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[error] %s:%s: %v", e.Component, e.Type, e.Err)
}

const (
	eConError   = "ECONNERROR"
	eValidError = "EVALIDERROR"
	eConFailed  = "ECONNFAILED"
)
