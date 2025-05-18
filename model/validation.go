package model

import (
	"errors"
	"reflect"
	"strings"
)

// Validatable interface for models that support validation
type Validatable interface {
	Validate() error
}

// Validation rules
const (
	Required = "required"
	Email    = "email"
	Min      = "min"
	Max      = "max"
)

// Rule defines a validation rule
type Rule struct {
	Name  string
	Value interface{}
}

// Field defines a validated field
type Field struct {
	Name  string
	Value interface{}
	Rules []Rule
}

// ValidateDocument validates a document before save
func ValidateDocument(doc interface{}) error {
	if v, ok := doc.(Validatable); ok {
		return v.Validate()
	}
	return nil
}

// ValidateFields validates fields against rules
func ValidateFields(fields []Field) error {
	var errs []string
	
	for _, field := range fields {
		for _, rule := range field.Rules {
			switch rule.Name {
			case Required:
				if isEmpty(field.Value) {
					errs = append(errs, field.Name+" is required")
				}
			case Email:
				// Email validation logic
			case Min:
				// Min length/value validation
			case Max:
				// Max length/value validation
			}
		}
	}
	
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}
