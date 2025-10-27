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

// Float converts an interface to a floating-point type, ignoring any conversion errors.
// It returns the zero value of the target type if conversion fails.
// E must be a floating-point type (float32 or float64).
//
// Example:
//
//	result := Float[float64]("3.14") // returns 3.14
//	result := Float[float32](true) // returns 1.0
//	result := Float[float64](42) // returns 42.0
func Float[E constraints.Float](o any) E {
	v, _ := FloatE[E](o)
	return v
}

// FloatE converts an interface to a floating-point type, returning both the converted value and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// E must be a floating-point type (float32 or float64).
//
// Example:
//
//	result, err := FloatE[float64]("3.14") // returns 3.14, nil
//	result, err := FloatE[float64]("invalid") // returns 0.0, error
func FloatE[E constraints.Float](o any) (E, error) {
	return floatE[E](o)
}

// FloatS converts an interface to a floating-point slice type, ignoring any conversion errors.
// It's designed for converting slice-like data structures to floating-point slices.
// S is a slice type with elements of floating-point type.
// E must be a floating-point type (float32 or float64).
//
// Example:
//
//	result := FloatS[[]float64]([]string{"1.1", "2.2", "3.3"}) // returns []float64{1.1, 2.2, 3.3}
func FloatS[S ~[]E, E constraints.Float](o any) S {
	v, _ := FloatSE[S](o)
	return v
}

// FloatSE converts an interface to a floating-point slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors for slice data explicitly.
// S is a slice type with elements of floating-point type.
// E must be a floating-point type (float32 or float64).
//
// Example:
//
//	result, err := FloatSE[[]float64]([]string{"1.1", "2.2"}) // returns []float64{1.1, 2.2}, nil
//	result, err := FloatSE[[]float64]([]string{"1.1", "invalid"}) // returns nil, error
func FloatSE[S ~[]E, E constraints.Float](o any) (S, error) {
	return toSliceE[S](o, floatE[E])
}

// floatE is the core implementation of floating-point conversion with error handling.
// It uses a fast path approach for common types and falls back to reflection for complex types.
// E must be a floating-point type (float32 or float64).
func floatE[E constraints.Float](o any) (E, error) {
	var zero E
	// Handle nil input by returning zero value
	if o == nil {
		return zero, nil
	}

	// Fast path: direct type assertions for common types
	switch f := o.(type) {
	// Boolean conversion: true becomes 1.0, false becomes 0.0
	case bool:
		if f {
			return 1, nil
		}
		return zero, nil

	// Native floating-point types
	case float64:
		return E(f), nil
	case float32:
		return E(f), nil

	// Integer types: direct conversion to floating-point
	case int:
		return E(f), nil
	case int64:
		return E(f), nil
	case int32:
		return E(f), nil
	case int16:
		return E(f), nil
	case int8:
		return E(f), nil
	case uint:
		return E(f), nil
	case uint64:
		return E(f), nil
	case uint32:
		return E(f), nil
	case uint16:
		return E(f), nil
	case uint8:
		return E(f), nil

	// String conversion using strconv.ParseFloat
	case string:
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Byte slice conversion by converting to string first
	case []byte:
		v, err := strconv.ParseFloat(string(f), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// JSON number support
	case json.Number:
		v, err := f.Float64()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Time types that can be converted to numeric values
	case time.Weekday:
		return E(f), nil
	case time.Month:
		return E(f), nil
	case time.Duration:
		return E(f), nil

	// Database driver.Valuer interface support
	case driver.Valuer:
		v, err := f.Value()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		r, err := floatE[E](v)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return r, nil

	// Protobuf duration type support
	case *durationpb.Duration:
		return E(f.AsDuration()), nil

	// Protobuf wrapper types support
	case *wrapperspb.BoolValue:
		if f.GetValue() {
			return 1, nil
		}
		return zero, nil
	case *wrapperspb.DoubleValue:
		return E(f.GetValue()), nil
	case *wrapperspb.FloatValue:
		return E(f.GetValue()), nil
	case *wrapperspb.Int64Value:
		return E(f.GetValue()), nil
	case *wrapperspb.Int32Value:
		return E(f.GetValue()), nil
	case *wrapperspb.UInt64Value:
		return E(f.GetValue()), nil
	case *wrapperspb.UInt32Value:
		return E(f.GetValue()), nil
	case *wrapperspb.StringValue:
		v, err := strconv.ParseFloat(f.GetValue(), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case *wrapperspb.BytesValue:
		v, err := strconv.ParseFloat(string(f.GetValue()), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Stringer interface support for custom types that can be represented as strings
	case fmt.Stringer:
		v, err := strconv.ParseFloat(f.String(), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Default case: use reflection-based conversion for complex types
	default:
		// slow path
		return floatVE[E](o)
	}
}

// floatVE is the reflection-based (slow path) implementation for floating-point conversion.
// It's used when fast path type assertions fail and more complex type analysis is needed.
// E must be a floating-point type (float32 or float64).
func floatVE[E constraints.Float](o any) (E, error) {
	var zero E
	// Get the underlying value, dereferencing pointers if necessary
	v := indirectValue(reflect.ValueOf(o))

	// Handle different reflection kinds
	switch v.Kind() {
	// Boolean conversion: true becomes 1.0, false becomes 0.0
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return zero, nil

	// Integer types: direct conversion to floating-point
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return E(v.Int()), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return E(v.Uint()), nil

	// Floating-point types: direct conversion
	case reflect.Float64, reflect.Float32:
		return E(v.Float()), nil

	// String conversion using strconv.ParseFloat
	case reflect.String:
		f, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(f), nil

	// Byte slice conversion (must be []byte)
	case reflect.Slice:
		// Ensure it's a byte slice
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		f, err := strconv.ParseFloat(string(v.Bytes()), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(f), nil

	// Unsupported types
	default:
		return failedCastValue[E](o)
	}
}
