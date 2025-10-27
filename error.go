package gonv

import "fmt"

// Error message templates for failed type conversions
var (
	// failedCast is the error message template for conversion failures without underlying error
	// Format: "gonv: failed to cast 'value' of type OriginalType to TargetType"
	failedCast = "gonv: failed to cast %#v of type %T to %T"

	// failedCastErr is the error message template for conversion failures with underlying error
	// Format: "gonv: failed to cast 'value' of type OriginalType to TargetType, underlying error"
	failedCastErr = failedCast + ", %w"
)

// failedCastValue creates a zero value of type E and returns it with a formatted error message.
// Used when a type conversion fails without an underlying error.
//
// Example:
//
//	var result, err = failedCastValue[int]("not a number")
//	// result = 0, err = "gonv: failed to cast "not a number" of type string to int"
func failedCastValue[E any](o any) (E, error) {
	var zero E
	return zero, fmt.Errorf(failedCast, o, o, zero)
}

// failedCastErrValue creates a zero value of type E and returns it with a formatted error message
// that includes the underlying error. Used when a type conversion fails with an underlying error.
//
// Example:
//
//	var result, err = failedCastErrValue[int]("not a number", strconv.ErrSyntax)
//	// result = 0, err = "gonv: failed to cast "not a number" of type string to int, strconv.ErrSyntax"
func failedCastErrValue[E any](o any, err error) (E, error) {
	var zero E
	return zero, fmt.Errorf(failedCastErr, o, o, zero, err)
}
