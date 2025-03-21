package data

import (
	"fmt"
	"strconv"
)

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
