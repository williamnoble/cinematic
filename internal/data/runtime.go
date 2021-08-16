package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Runtime is a custom integer type representing the total runtime
// for a given movie
type Runtime int32

// ErrInvalidRuntimeFormat => runtime formatting error
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// MarshalJSON ensures the datatype Runtype int32 is Marshalled correctly.
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJsonValue := strconv.Quote(jsonValue)
	return []byte(quotedJsonValue), nil
}

// UnmarshalJSON converts custom Runtime string to an int32
// Example. Movie.Runtime = "103 Mins"
// 1. Unquote the string "103 Mins", 2. Split by Whitespace
// yielding []string{"103", "mins"}. Part[0]string="103"
// 3. Convert string to an int32 value.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// return custom type
	*r = Runtime(i)

	return nil
}
