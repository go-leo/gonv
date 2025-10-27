package gonv

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Uint converts an interface to an unsigned integer type, ignoring any conversion errors.
// It returns the zero value of the target type if conversion fails.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
//
// Example:
//
//	result := Uint[uint64]("42") // returns 42
//	result := Uint[uint](true) // returns 1
//	result := Uint[uint](3.14) // returns 3
func Uint[E constraints.Unsigned](o any) E {
	v, _ := UintE[E](o)
	return v
}

// UintE converts an interface to an unsigned integer type, returning both the converted value and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
//
// Example:
//
//	result, err := UintE[uint64]("42") // returns 42, nil
//	result, err := UintE[uint64]("-1") // returns 0, error (negative values are not allowed)
func UintE[E constraints.Unsigned](o any) (E, error) {
	return uintE[E](o)
}

// UintS converts an interface to an unsigned integer slice type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
// S is a slice type with elements of unsigned integer type.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
//
// Example:
//
//	result := UintS[[]uint64, uint64]([]string{"1", "2", "3"}) // returns []uint64{1, 2, 3}
func UintS[S ~[]E, E constraints.Unsigned](o any) S {
	v, _ := UintSE[S](o)
	return v
}

// UintSE converts an interface to an unsigned integer slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors for slice data explicitly.
// S is a slice type with elements of unsigned integer type.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
//
// Example:
//
//	result, err := UintSE[[]uint64, uint64]([]string{"1", "2"}) // returns []uint64{1, 2}, nil
//	result, err := UintSE[[]uint64, uint64]([]string{"1", "-1"}) // returns nil, error (negative values are not allowed)
func UintSE[S ~[]E, E constraints.Unsigned](o any) (S, error) {
	return toSliceE[S](o, uintE[E])
}

// uintE is the core implementation of unsigned integer conversion with error handling.
// It uses a fast path approach for common types and falls back to reflection for complex types.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
// Negative values are not allowed and will result in an error.
func uintE[E constraints.Unsigned](o any) (E, error) {
	var zero E
	// Handle nil input by returning zero value
	if o == nil {
		return zero, nil
	}

	// Fast path: direct type assertions for common types
	switch u := o.(type) {
	// Boolean conversion: true becomes 1, false becomes 0
	case bool:
		if u {
			return 1, nil
		}
		return zero, nil

	// Signed integer types: check for negative values
	case int:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case int64:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case int32:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case int16:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case int8:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil

	// Unsigned integer types: direct conversion
	case uint:
		return E(u), nil
	case uint64:
		return E(u), nil
	case uint32:
		return E(u), nil
	case uint16:
		return E(u), nil
	case uint8:
		return E(u), nil

	// Floating-point types: check for negative values
	case float64:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case float32:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil

	// String conversion using strconv.ParseUint with trimZeroDecimal
	case string:
		v, err := strconv.ParseUint(trimZeroDecimal(u), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Byte slice conversion by converting to string first
	case []byte:
		v, err := strconv.ParseUint(trimZeroDecimal(string(u)), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Time types that can be converted to numeric values
	case time.Duration:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case time.Weekday:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case time.Month:
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil

	// JSON number support
	case json.Number:
		v, err := u.Int64()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), err

	// Protobuf duration type support: convert to duration then check for negative values
	case *durationpb.Duration:
		v := u.AsDuration()
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), nil

	// Protobuf wrapper types support
	case *wrapperspb.BoolValue:
		if u.GetValue() {
			return 1, nil
		}
		return zero, nil
	case *wrapperspb.Int64Value:
		v := u.GetValue()
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), nil
	case *wrapperspb.Int32Value:
		v := u.GetValue()
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), nil
	case *wrapperspb.UInt64Value:
		return E(u.GetValue()), nil
	case *wrapperspb.UInt32Value:
		return E(u.GetValue()), nil
	case *wrapperspb.DoubleValue:
		v := u.GetValue()
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), nil
	case *wrapperspb.FloatValue:
		v := u.GetValue()
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), nil
	case *wrapperspb.StringValue:
		v, err := strconv.ParseUint(trimZeroDecimal(u.GetValue()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case *wrapperspb.BytesValue:
		v, err := strconv.ParseUint(trimZeroDecimal(string(u.GetValue())), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Database driver.Valuer interface support
	case driver.Valuer:
		v, err := u.Value()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		r, err := uintE[E](v)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return r, nil

	// Stringer interface support for custom types that can be represented as strings
	case fmt.Stringer:
		v, err := strconv.ParseUint(trimZeroDecimal(string(u.String())), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Default case: use reflection-based conversion for complex types
	default:
		return toUnsignedValueE[E](o)
	}
}

// toUnsignedValueE is the reflection-based (slow path) implementation for unsigned integer conversion.
// It's used when fast path type assertions fail and more complex type analysis is needed.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
// Negative values are not allowed and will result in an error.
func toUnsignedValueE[E constraints.Unsigned](o any) (E, error) {
	// Get the underlying value, dereferencing pointers if necessary
	v := indirectValue(reflect.ValueOf(o))
	var zero E

	// Handle different reflection kinds
	switch v.Kind() {
	// Boolean conversion: true becomes 1, false becomes 0
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return zero, nil

	// Signed integer types: check for negative values
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		u := v.Int()
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil

	// Unsigned integer types: direct conversion
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return E(v.Uint()), nil

	// Floating-point types: check for negative values
	case reflect.Float64, reflect.Float32:
		u := v.Float()
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil

	// String conversion using strconv.ParseUint with trimZeroDecimal
	case reflect.String:
		u, err := strconv.ParseUint(trimZeroDecimal(v.String()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(u), nil

	// Byte slice conversion (must be []byte)
	case reflect.Slice:
		// Ensure it's a byte slice
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		u, err := strconv.ParseUint(trimZeroDecimal(string(v.Bytes())), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(u), nil

	// Unsupported types
	default:
		return failedCastValue[E](o)
	}
}
