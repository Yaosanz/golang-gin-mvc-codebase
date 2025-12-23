package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppValidator struct {
	validate *validator.Validate
}

func NewAppValidator(db *gorm.DB) *AppValidator {
	v := validator.New()

	// register custom validator
	_ = v.RegisterValidation("unique", unique(db))
	_ = v.RegisterValidation("exists", exists(db))
	_ = v.RegisterValidation("required_top_field_if", requiredTopFieldIf)
	_ = v.RegisterValidation("required_if_not_empty", requiredIfNotEmpty)

	return &AppValidator{validate: v}
}

func (v *AppValidator) ValidateStruct(i interface{}) error {
	err := v.validate.Struct(i)
	if err != nil {
		return ParseValidationErrors(err, i)
	}
	return nil
}

// camelCase is helper function to convert snake_case to CamelCase.
func camelCase(input string) string {
	isToUpper := false
	result := ""

	for i, v := range input {
		if i == 0 {
			// Capitalize the first letter
			result += strings.ToUpper(string(v))
		} else if v == '_' {
			// Set flag to capitalize the next character
			isToUpper = true
		} else {
			if isToUpper {
				// Capitalize if the previous character was an underscore
				result += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				result += string(v)
			}
		}
	}

	// Ensure 'ID' is always uppercase
	if strings.HasSuffix(result, "Id") {
		result = result[:len(result)-2] + "ID"
	}

	return result
}

// getJSONTag Helper function to get the json tag from the struct field
func getJSONTag(structType reflect.Type, field string) string {
	fieldStruct, found := structType.FieldByName(field)
	if !found {
		return field // If no json tag is found, return the original field name
	}

	jsonTag := fieldStruct.Tag.Get("json")
	if jsonTag == "" {
		return field // If no json tag, return the original field name
	}

	// Handle cases where the json tag has options like "omitempty"
	return strings.Split(jsonTag, ",")[0]
}
