// Package gonv provides type conversion utilities for Go applications.
// It offers safe and flexible casting between different data types with generic support.
// This file contains functions for converting values to slice types.
package gonv

import (
	"reflect"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// AnySlice casts an interface to a []any type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
//
// Example:
//
//	result := AnySlice([]string{"a", "b", "c"}) // returns []interface{}{"a", "b", "c"}
//	result := AnySlice([3]string{"a", "b", "c"}) // returns []interface{}{"a", "b", "c"}
func AnySlice(o any) []any {
	v, _ := AnySliceE(o)
	return v
}

// AnySliceE casts an interface to a []any type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
//
// Example:
//
//	result, err := AnySliceE([]string{"a", "b", "c"}) // returns []interface{}{"a", "b", "c"}, nil
//	result, err := AnySliceE("not a slice") // returns nil, error
func AnySliceE(o any) ([]any, error) {
	return toSliceE[[]any](o, func(o any) (any, error) { return o, nil })
}

// StringSlice casts an interface to a string slice type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
// S is a slice type with elements of string type.
// E must be a string type.
//
// Example:
//
//	result := StringSlice[[]string, string]([]string{"a", "b", "c"}) // returns []string{"a", "b", "c"}
func StringSlice[S ~[]E, E ~string](o any) S {
	v, _ := StringSliceE[S](o)
	return v
}

// StringSliceE casts an interface to a string slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// S is a slice type with elements of string type.
// E must be a string type.
//
// Example:
//
//	result, err := StringSliceE[[]string, string]([]string{"a", "b", "c"}) // returns []string{"a", "b", "c"}, nil
//	result, err := StringSliceE[[]string, string]("not a slice") // returns nil, error
func StringSliceE[S ~[]E, E ~string](o any) (S, error) {
	return toSliceE[S](o, StringE[E])
}

// BoolSlice casts an interface to a boolean slice type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
// S is a slice type with elements of boolean type.
// E must be a boolean type.
//
// Example:
//
//	result := BoolSlice[[]bool, bool]([]bool{true, false, true}) // returns []bool{true, false, true}
func BoolSlice[S ~[]E, E ~bool](o any) S {
	v, _ := BoolSliceE[S](o)
	return v
}

// BoolSliceE casts an interface to a boolean slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// S is a slice type with elements of boolean type.
// E must be a boolean type.
//
// Example:
//
//	result, err := BoolSliceE[[]bool, bool]([]bool{true, false, true}) // returns []bool{true, false, true}, nil
//	result, err := BoolSliceE[[]bool, bool]("not a slice") // returns nil, error
func BoolSliceE[S ~[]E, E ~bool](o any) (S, error) {
	return toSliceE[S](o, BoolE[E])
}

// FloatSlice casts an interface to a floating-point slice type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
// S is a slice type with elements of floating-point type.
// E must be a floating-point type (float32 or float64).
//
// Example:
//
//	result := FloatSlice[[]float64, float64]([]float64{1.1, 2.2, 3.3}) // returns []float64{1.1, 2.2, 3.3}
func FloatSlice[S ~[]E, E constraints.Float](o any) S {
	v, _ := FloatSliceE[S](o)
	return v
}

// FloatSliceE casts an interface to a floating-point slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// S is a slice type with elements of floating-point type.
// E must be a floating-point type (float32 or float64).
//
// Example:
//
//	result, err := FloatSliceE[[]float64, float64]([]float64{1.1, 2.2, 3.3}) // returns []float64{1.1, 2.2, 3.3}, nil
//	result, err := FloatSliceE[[]float64, float64]("not a slice") // returns nil, error
func FloatSliceE[S ~[]E, E constraints.Float](o any) (S, error) {
	return toSliceE[S](o, FloatE[E])
}

// IntSlice casts an interface to a signed integer slice type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
// S is a slice type with elements of signed integer type.
// E must be a signed integer type (int, int8, int16, int32, int64).
//
// Example:
//
//	result := IntSlice[[]int64, int64]([]int64{1, 2, 3}) // returns []int64{1, 2, 3}
func IntSlice[S ~[]E, E constraints.Signed](o any) S {
	v, _ := IntSliceE[S](o)
	return v
}

// IntSliceE casts an interface to a signed integer slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// S is a slice type with elements of signed integer type.
// E must be a signed integer type (int, int8, int16, int32, int64).
//
// Example:
//
//	result, err := IntSliceE[[]int64, int64]([]int64{1, 2, 3}) // returns []int64{1, 2, 3}, nil
//	result, err := IntSliceE[[]int64, int64]("not a slice") // returns nil, error
func IntSliceE[S ~[]E, E constraints.Signed](o any) (S, error) {
	return toSliceE[S](o, IntE[E])
}

// UintSlice casts an interface to an unsigned integer slice type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
// S is a slice type with elements of unsigned integer type.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
//
// Example:
//
//	result := UintSlice[[]uint64, uint64]([]uint64{1, 2, 3}) // returns []uint64{1, 2, 3}
func UintSlice[S ~[]E, E constraints.Unsigned](o any) S {
	v, _ := UintSliceE[S](o)
	return v
}

// UintSliceE casts an interface to an unsigned integer slice type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// S is a slice type with elements of unsigned integer type.
// E must be an unsigned integer type (uint, uint8, uint16, uint32, uint64).
//
// Example:
//
//	result, err := UintSliceE[[]uint64, uint64]([]uint64{1, 2, 3}) // returns []uint64{1, 2, 3}, nil
//	result, err := UintSliceE[[]uint64, uint64]("not a slice") // returns nil, error
func UintSliceE[S ~[]E, E constraints.Unsigned](o any) (S, error) {
	return toSliceE[S](o, UintE[E])
}

// Slice is a generic function that casts an interface to a slice type, returning both the converted slice and any error encountered.
// S is a slice type with elements of type E.
// E is the element type of the slice.
// to is a function that converts individual elements.
//
// Example:
//
//	result, err := Slice[[]int, int]([]string{"1", "2", "3"}, IntE[int])
//	// returns []int{1, 2, 3}, nil
func Slice[S ~[]E, E any](o any, to func(o any) (E, error)) (S, error) {
	return SliceE[S](o, to)
}

// SliceE is a generic function that casts an interface to a slice type, returning both the converted slice and any error encountered.
// S is a slice type with elements of type E.
// E is the element type of the slice.
// to is a function that converts individual elements.
//
// Example:
//
//	result, err := SliceE[[]int, int]([]string{"1", "2", "3"}, IntE[int])
//	// returns []int{1, 2, 3}, nil
func SliceE[S ~[]E, E any](o any, to func(o any) (E, error)) (S, error) {
	return toSliceE[S, E](o, to)
}

// toSliceE is the core implementation of slice conversion with error handling.
// It uses type assertion for direct slice types and reflection for array/slice conversion.
// S is a slice type with elements of type E.
// E is the element type of the slice.
// to is a function that converts individual elements.
func toSliceE[S ~[]E, E any](o any, to func(o any) (E, error)) (S, error) {
	var zero S
	// Handle nil input by returning zero value
	if o == nil {
		return zero, nil
	}

	// Handle direct type match by cloning the slice
	if v, ok := o.(S); ok {
		return slices.Clone(v), nil
	}

	// Check if input is a slice or array type
	kind := reflect.TypeOf(o).Kind()
	switch kind {
	// Handle slice and array types by converting each element
	case reflect.Slice, reflect.Array:
		value := reflect.ValueOf(o)
		res := make(S, value.Len())
		for i := 0; i < value.Len(); i++ {
			val, err := to(value.Index(i).Interface())
			if err != nil {
				return zero, err
			}
			res[i] = val
		}
		return res, nil
	// Handle unsupported types
	default:
		return failedCastValue[S](o)
	}
}
