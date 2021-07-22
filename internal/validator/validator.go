package validator

import "regexp"

var (

	//	EmailRX is a regular expression which confirms the validity of an email input.
	//EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\. [a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator provides an error map for mapping error values encountered
type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if one or more errors are present
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error msg to the map as no entry already exists for key
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check (the main function) if error present, call AddError Function
// v.Check(input.Title != "", "title", "must be provided"
func (v *Validator) Check(ok bool, key string, message string) {
	if !ok { //if false
		v.AddError(key, message)
	}
}

// In checks the given string is contained in a list
func In(value string, list ...string) bool {
	for item := range list {
		if value == list[item] {
			return true
		}
	}
	return false
}

// Matches uses Regular Expressions to check a value
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value) //282
}

// Unique checks all values in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
