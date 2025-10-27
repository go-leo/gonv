package gonv

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Bool casts an interface to a bool type, ignoring any conversion errors.
// It returns the zero value of type E if conversion fails.
//
// Example:
//
//	result := Bool[bool]("true") // returns true
//	result := Bool[bool]("false") // returns false
//	result := Bool[bool](1) // returns true
//	result := Bool[bool](0) // returns false
func Bool[E ~bool](o any) E {
	v, _ := BoolE[E](o)
	return v
}

// BoolE casts an interface to a bool type, returning both the converted value and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
//
// Example:
//
//	result, err := BoolE[bool]("true") // returns true, nil
//	result, err := BoolE[bool]("invalid") // returns false, error
func BoolE[E ~bool](o any) (E, error) {
	return boolE[E](o)
}

// BoolS casts an interface to a []bool type, ignoring any conversion errors.
// It's designed for converting slice-like data structures to boolean slices.
//
// Example:
//
//	result := BoolS[[]bool]([]string{"true", "false", "1", "0"}) // returns []bool{true, false, true, false}
func BoolS[S ~[]E, E ~bool](o any) S {
	v, _ := BoolSE[S](o)
	return v
}

// BoolSE casts an interface to a []bool type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors for slice data explicitly.
//
// Example:
//
//	result, err := BoolSE[[]bool]([]string{"true", "false"}) // returns []bool{true, false}, nil
//	result, err := BoolSE[[]bool]([]string{"true", "invalid"}) // returns []bool{true, false}, error
func BoolSE[S ~[]E, E ~bool](o any) (S, error) {
	return toSliceE[S](o, boolE[E])
}

// boolE is the core implementation of boolean conversion with error handling.
// It uses a fast path approach for common types and falls back to reflection for complex types.
func boolE[E ~bool](o any) (E, error) {
	// Handle nil input by returning the zero value of type E
	if o == nil {
		var zero E
		return zero, nil
	}

	// Fast path: direct type assertions for common types
	switch b := o.(type) {
	// Native boolean type
	case bool:
		return E(b), nil

	// String conversion using strconv.ParseBool
	case string:
		v, err := strconv.ParseBool(b)
		if err != nil {
			return failedCastErrValue[E](b, err)
		}
		return E(v), err

	// Byte slice conversion by converting to string first
	case []byte:
		v, err := strconv.ParseBool(string(b))
		if err != nil {
			return failedCastErrValue[E](b, err)
		}
		return E(v), err

	// Protobuf wrapper types support
	case *wrapperspb.BoolValue:
		return E(b.GetValue()), nil
	case *wrapperspb.StringValue:
		v, err := strconv.ParseBool(b.GetValue())
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), err
	case *wrapperspb.BytesValue:
		v, err := strconv.ParseBool(string(b.GetValue()))
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), err

	// Numeric types: convert to float64 first, then check if non-zero
	case
		float64, float32,
		int, int64, int32, int16, int8,
		uint, uint64, uint32, uint16, uint8,
		json.Number,
		*wrapperspb.DoubleValue, *wrapperspb.FloatValue,
		*wrapperspb.Int64Value, *wrapperspb.Int32Value,
		*wrapperspb.UInt64Value, *wrapperspb.UInt32Value:
		n, err := FloatE[float64](o)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		// Non-zero numeric values are treated as true, zero as false
		return n != 0, nil

	// Database driver.Valuer interface support
	case driver.Valuer:
		v, err := b.Value()
		if err != nil {
			return failedCastErrValue[E](b, err)
		}
		r, err := boolE[E](v)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return r, nil

	// Stringer interface support for custom types that can be represented as strings
	case fmt.Stringer:
		v, err := strconv.ParseBool(b.String())
		if err != nil {
			return failedCastErrValue[E](b, err)
		}
		return E(v), err

	// Default case: use reflection-based conversion for complex types
	default:
		// slow path
		return boolVE[E](o)
	}
}

// boolVE is the reflection-based (slow path) implementation for boolean conversion.
// It's used when fast path type assertions fail and more complex type analysis is needed.
func boolVE[E ~bool](o any) (E, error) {
	// Get the underlying value, dereferencing pointers if necessary
	v := indirectValue(reflect.ValueOf(o))

	// Handle different reflection kinds
	switch v.Kind() {
	// Native boolean type
	case reflect.Bool:
		return E(v.Bool()), nil

	// Integer types: non-zero values are true, zero is false
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return v.Int() != 0, nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return v.Uint() != 0, nil

	// Floating point types: non-zero values are true, zero is false
	case reflect.Float64, reflect.Float32:
		return v.Float() != 0, nil

	// String conversion using strconv.ParseBool
	case reflect.String:
		b, err := strconv.ParseBool(v.String())
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(b), err

	// Byte slice conversion (must be []byte)
	case reflect.Slice:
		// Ensure it's a byte slice
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		b, err := strconv.ParseBool(string(v.Bytes()))
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(b), err

	// Unsupported types
	default:
		return failedCastValue[E](o)
	}
}
