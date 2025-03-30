package validator

import (
	"regexp"
	"slices"
)

// declare regex for sanity checking the validity of email addresses
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// define validator struct that contains a map of validation errors
type Validator struct {
	Errors map[string]string
}

// New is a helper which creates a new Validator instance with an empty errors map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// the Valid method returns true if the errors map does not contain any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError method adds an error to the map (so long as no entry exists for the given key)
func (v *Validator) AddError(key, message string) {
	// check if the value exists for the key
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// the Check method adds an error to the map, only if the validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Generic function which returns true if a specific value is in a list of permitted values
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// the Matches method returns true if a string value matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Generic function which returns true if all the values in a slice are unique
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
