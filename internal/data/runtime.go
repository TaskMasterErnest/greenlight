package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// define error to be returned by UnmarshalJSON()
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// declare a custom Runtime type, same as in the Movie struct
type Runtime int32

// implement a MarshalJSON() method on the Runtime type, to satisfy the json.Marshaler interface
func (r Runtime) MarshalJSON() ([]byte, error) {
	// generate a string containing the movie runtime value in the required format
	jsonValue := fmt.Sprintf("%d mins", r)

	// wrap it in double quotes with the strconv.Quote() function
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// implement UnmarshalJSON() method on the Runtime type so that it satisfies the json.Unmarshaler interface
// it must be a pointer so that we modify the underlying value itself instead of a copy
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// the expected incoming input for the runtime is a string ("<runtime> mins"), so we unquote it
	unquotedJSONvalue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// we split the string to isolate both parts
	parts := strings.Split(unquotedJSONvalue, " ")

	// perform sanity checks on both parts
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// parse the part that contains the number into an int32 type
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// convert the int32 into a Runtime type and assign this to a receiver
	// dereference the receiver in order to set the underlying value of the pointer
	*r = Runtime(i)

	return nil
}
