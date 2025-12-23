package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Error validationError is a custom error type that holds field-specific errors.
type Error struct {
	Errors []ErrorFormat `json:"errors"`
}

func (v *Error) Error() string {
	return "Validation failed"
}

type ErrorFormat struct {
	Tag     string `json:"tag"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ParseValidationErrors processes validation errors and returns a map with field names and custom error messages.
func ParseValidationErrors(err error, input interface{}) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		format := make([]ErrorFormat, 0)

		for _, fieldErr := range validationErrors {
			customField := buildCustomFieldName(input, fieldErr)
			message := getErrorMessage(fieldErr, customField.titleTag, "id")
			format = append(format, ErrorFormat{
				Tag:     fieldErr.Tag(),
				Field:   customField.namespace,
				Message: message,
			})
		}

		return &Error{Errors: format}
	}

	return nil
}

type fieldFormat struct {
	namespace string
	titleTag  string
	jsonTag   string
}

// buildFullFieldNamespace Build full field name for nested structs and slices
func buildFullFieldNamespace(err validator.FieldError) string {
	// Get the struct field names and join them using '.'
	ns := err.Namespace()                 // e.g. "User.Address.Zip" or "User.Cars[0].Brand"
	nsParts := strings.Split(ns, ".")[1:] // Remove the root struct name (e.g. "User")
	return strings.Join(nsParts, ".")
}

// buildCustomFieldName Build full field name for nested structs and slices
func buildCustomFieldName(input interface{}, err validator.FieldError) *fieldFormat {
	parts := strings.Split(err.Namespace(), ".")
	var fieldNamespace []string
	var fieldTitleTag string
	var fieldJsonTag string

	// Dereference the input if it's a pointer
	currentInputValue := reflect.ValueOf(input)
	currentInputType := reflect.TypeOf(input)
	if currentInputType.Kind() == reflect.Ptr {
		currentInputValue = currentInputValue.Elem()
		currentInputType = currentInputType.Elem()
	}

	for i := 1; i < len(parts); i++ { // Skip the root struct name
		fieldName := parts[i]

		// Handle array index (e.g., Data[0])
		if strings.Contains(fieldName, "[") {
			getName := getArrayName(fieldName)

			// Lookup field by name in the current struct
			field, found := currentInputType.FieldByName(getName)
			if !found {
				continue // Ignore if the field isn't found
			}

			// Get the json tag or fallback to the field name
			jsonTag := getFieldTag(field, "json", field.Name)
			fieldJsonTag = jsonTag

			// Get the title tag or fallback to the field name
			titleTag := getFieldTag(field, "title", "")
			fieldTitleTag = titleTag

			// current field value
			fieldValue := currentInputValue.FieldByName(getName)
			arrayIndex := getArrayIndex(fieldName)
			if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
				// Get the element at the array index
				currentInputValue = fieldValue.Index(parseArrayIndex(arrayIndex))
				currentInputType = currentInputValue.Type() // Move to the element's type
			}

			// set array name as the parent namespace
			fieldNamespace = append(fieldNamespace, jsonTag)
			// append the array index to last added namespace
			fieldNamespace[len(fieldNamespace)-1] += fmt.Sprintf(".%s", arrayIndex)
			continue
		}

		// Lookup field by name in the current struct
		field, found := currentInputType.FieldByName(fieldName)
		if !found {
			continue // Ignore if the field isn't found
		}

		// Get the json tag or fallback to the field name
		jsonTag := getFieldTag(field, "json", field.Name)
		fieldJsonTag = jsonTag // set fieldJsonTag value

		// Get the title tag or fallback to the json name
		titleTag := getFieldTag(field, "title", "")
		fieldTitleTag = titleTag // set fieldTitleTag value

		// set fieldNamespace value
		fieldNamespace = append(fieldNamespace, jsonTag)

		// Move to the next nested field (whether it's struct or slice)
		currentInputValue = currentInputValue.FieldByName(fieldName)
		currentInputType = field.Type

		// Handle pointer dereferencing
		if currentInputType.Kind() == reflect.Ptr {
			currentInputValue = currentInputValue.Elem()
			currentInputType = currentInputType.Elem()
		}
	}

	return &fieldFormat{
		namespace: strings.Join(fieldNamespace, "."),
		titleTag:  fieldTitleTag,
		jsonTag:   fieldJsonTag,
	}
}

// getArrayName helper function to remove array index from a string like "Data[0]" to "Data"
func getArrayName(s string) string {
	// Find the position of the first '['
	index := strings.Index(s, "[")
	if index != -1 {
		// Return the substring before the '['
		return s[:index]
	}
	// If no '[' is found, return the original string (could be useful for cases without index)
	return s
}

// getArrayIndex Helper function to extract the array index from a string like "Data[0]" to "0"
func getArrayIndex(fieldName string) string {
	start := strings.Index(fieldName, "[") + 1
	end := strings.Index(fieldName, "]")
	return fieldName[start:end]
}

// parseArrayIndex Helper function to parse the array index to an integer
func parseArrayIndex(index string) int {
	i, _ := strconv.Atoi(index)
	return i
}

// getFieldTag Helper function to extract field tag value
func getFieldTag(structField reflect.StructField, tagName string, fallback string) string {
	tag := structField.Tag.Get(tagName)
	if tag == "" || tag == "-" {
		tag = fallback
	}
	return tag
}

// getErrorMessage validation custom message
func getErrorMessage(fieldErr validator.FieldError, title string, lang string) string {
	messages := map[string]map[string]string{
		"en": {
			"required":              fmt.Sprintf("is required"),
			"required_if":           fmt.Sprintf("is required"),
			"required_top_field_if": fmt.Sprintf("is required"),
			"required_if_not_empty": fmt.Sprintf("is required"),
			"email":                 fmt.Sprintf("must be a valid email"),
			"gte":                   fmt.Sprintf("must be greater than or equal to the required value"),
			"len":                   fmt.Sprintf("must be of length %s", fieldErr.Param()),
			"datetime":              fmt.Sprintf("not valid"),
		},
		"id": { // Indonesian translations
			"required":              strings.TrimSpace(fmt.Sprintf("%s harus diisi", title)),
			"required_if":           strings.TrimSpace(fmt.Sprintf("%s harus diisi", title)),
			"required_top_field_if": strings.TrimSpace(fmt.Sprintf("%s harus diisi", title)),
			"required_if_not_empty": strings.TrimSpace(fmt.Sprintf("%s harus diisi", title)),
			"email":                 strings.TrimSpace(fmt.Sprintf("%s harus berupa email yang valid", title)),
			"gte":                   strings.TrimSpace(fmt.Sprintf("%s harus lebih besar atau sama dengan %s", title, fieldErr.Param())),
			"len":                   strings.TrimSpace(fmt.Sprintf("%s harus memiliki panjang %s", title, fieldErr.Param())),
			"min":                   strings.TrimSpace(fmt.Sprintf("%s minimal %s", title, fieldErr.Param())),
			"unique":                strings.TrimSpace(fmt.Sprintf("%s sudah di gunakan", title)),
			"exists":                strings.TrimSpace(fmt.Sprintf("%s tidak valid", title)),
			"datetime":              strings.TrimSpace(fmt.Sprintf("%s tidak valid", title)),
		},
	}

	// Default message
	message := messages["id"][fieldErr.Tag()]

	// lang is set
	if lang != "" {
		if msg, ok := messages[lang][fieldErr.Tag()]; ok {
			message = msg
		}
	}

	// if no message fallback to system default message
	if message == "" {
		message = fmt.Sprintf("Field '%s' failed on the '%s' tag", fieldErr.Field(), fieldErr.Tag())
	}

	return message
}
