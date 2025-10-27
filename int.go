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
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Int converts an interface to a signed integer type, ignoring any conversion errors.
// It returns the zero value of the target type if conversion fails.
// E must be a signed integer type (int, int8, int16, int32, int64).
//
// Example:
//
//	result := Int[int64]("42") // returns 42
//	result := Int[int](true) // returns 1
//	result := Int[int](3.14) // returns 3
func Int[E constraints.Signed](o any) E {
	v, _ := IntE[E](o)
	return v
}

// IntE converts an interface to a signed integer type, returning both the converted value and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// E must be a signed integer type (int, int8, int16, int32, int64).
//
// Example:
//
//	result, err := IntE[int64]("42") // returns 42, nil
//	result, err := IntE[int64]("invalid") // returns 0, error
func IntE[E constraints.Signed](o any) (E, error) {
	return intE[E](o)
}

// IntS converts an interface to a signed integer slice type, ignoring any conversion errors.
// It's designed for converting slice-like data structures to signed integer slices.
// S is a slice type with elements of signed integer type.
// E must be a signed integer type (int, int8, int16, int32, int64).
//
// Example:
//
//	result := IntS[[]int64]([]string{"1", "2", "3"}) // returns []int64{1, 2, 3}
func IntS[S ~[]E, E constraints.Signed](o any) S {
	v, _ := IntSE[S](o)
	return v
}

// IntSE converts an interface to a signed integer slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors for slice data explicitly.
// S is a slice type with elements of signed integer type.
// E must be a signed integer type (int, int8, int16, int32, int64).
//
// Example:
//
//	result, err := IntSE[[]int64]([]string{"1", "2"}) // returns []int64{1, 2}, nil
//	result, err := IntSE[[]int64]([]string{"1", "invalid"}) // returns nil, error
func IntSE[S ~[]E, E constraints.Signed](o any) (S, error) {
	return toSliceE[S](o, IntE[E])
}

// intE is the core implementation of signed integer conversion with error handling.
// It uses a fast path approach for common types and falls back to reflection for complex types.
// E must be a signed integer type (int, int8, int16, int32, int64).
func intE[E constraints.Signed](o any) (E, error) {
	var zero E
	// Handle nil input by returning zero value
	if o == nil {
		return zero, nil
	}

	// Fast path: direct type assertions for common types
	switch s := o.(type) {
	// Boolean conversion: true becomes 1, false becomes 0
	case bool:
		if s {
			return 1, nil
		}
		return zero, nil

	// Floating-point types: direct conversion to integer (truncates decimal part)
	case float64:
		return E(s), nil
	case float32:
		return E(s), nil

	// Integer types: direct conversion
	case int:
		return E(s), nil
	case int64:
		return E(s), nil
	case int32:
		return E(s), nil
	case int16:
		return E(s), nil
	case int8:
		return E(s), nil
	case uint:
		return E(s), nil
	case uint64:
		return E(s), nil
	case uint32:
		return E(s), nil
	case uint16:
		return E(s), nil
	case uint8:
		return E(s), nil

	// String conversion using strconv.ParseInt with trimZeroDecimal
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Byte slice conversion by converting to string first
	case []byte:
		v, err := strconv.ParseInt(trimZeroDecimal(string(s)), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// JSON number support
	case json.Number:
		v, err := s.Int64()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Time types that can be converted to numeric values
	case time.Weekday:
		return E(s), nil
	case time.Month:
		return E(s), nil
	case time.Duration:
		return E(s), nil

	// Protobuf duration type support: convert to duration then to integer
	case *durationpb.Duration:
		return E(s.AsDuration()), nil

	// Protobuf timestamp type support: convert to milliseconds since Unix epoch
	case *timestamppb.Timestamp:
		return E(s.AsTime().UnixMilli()), nil

	// Protobuf wrapper types support
	case *wrapperspb.BoolValue:
		if s.GetValue() {
			return 1, nil
		}
		return zero, nil
	case *wrapperspb.DoubleValue:
		return E(s.GetValue()), nil
	case *wrapperspb.FloatValue:
		return E(s.GetValue()), nil
	case *wrapperspb.Int64Value:
		return E(s.GetValue()), nil
	case *wrapperspb.Int32Value:
		return E(s.GetValue()), nil
	case *wrapperspb.UInt64Value:
		return E(s.GetValue()), nil
	case *wrapperspb.UInt32Value:
		return E(s.GetValue()), nil
	case *wrapperspb.StringValue:
		i, err := strconv.ParseInt(trimZeroDecimal(s.GetValue()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(i), nil
	case *wrapperspb.BytesValue:
		i, err := strconv.ParseInt(trimZeroDecimal(string(s.GetValue())), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(i), nil

	// Database driver.Valuer interface support
	case driver.Valuer:
		v, err := s.Value()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		r, err := intE[E](v)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return r, nil

	// Stringer interface support for custom types that can be represented as strings
	case fmt.Stringer:
		v, err := strconv.ParseInt(trimZeroDecimal(s.String()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil

	// Default case: use reflection-based conversion for complex types
	default:
		// slow path
		return toSignedValueE[E](o)
	}
}

// toSignedValueE is the reflection-based (slow path) implementation for signed integer conversion.
// It's used when fast path type assertions fail and more complex type analysis is needed.
// E must be a signed integer type (int, int8, int16, int32, int64).
func toSignedValueE[E constraints.Signed](o any) (E, error) {
	var zero E
	// Get the underlying value, dereferencing pointers if necessary
	v := indirectValue(reflect.ValueOf(o))

	// Handle different reflection kinds
	switch v.Kind() {
	// Boolean conversion: true becomes 1, false becomes 0
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return zero, nil

	// Integer types: direct conversion
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return E(v.Int()), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return E(v.Uint()), nil

	// Floating-point types: conversion to integer (truncates decimal part)
	case reflect.Float64, reflect.Float32:
		return E(v.Float()), nil

	// String conversion using strconv.ParseInt with trimZeroDecimal
	case reflect.String:
		i, err := strconv.ParseInt(trimZeroDecimal(v.String()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(i), nil

	// Byte slice conversion (must be []byte)
	case reflect.Slice:
		// Ensure it's a byte slice
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		i, err := strconv.ParseInt(trimZeroDecimal(string(v.Bytes())), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(i), nil

	// Unsupported types
	default:
		return failedCastValue[E](o)
	}
}
