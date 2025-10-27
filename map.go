package gonv

import (
	"encoding/json"
	"reflect"

	"golang.org/x/exp/constraints"
)

// StringAnyMap casts an interface to a map[string]any type, ignoring any conversion errors.
// It returns an empty map if conversion fails.
// K must be a string type.
//
// Example:
//
//	result := StringAnyMap[string](map[string]interface{}{"key": "value"}) // returns map[string]interface{}{"key": "value"}
//	result := StringAnyMap[string](`{"key": "value"}`) // returns map[string]interface{}{"key": "value"}
func StringAnyMap[K ~string](o any) map[K]any {
	v, _ := StringAnyMapE[K](o)
	return v
}

// StringAnyMapE casts an interface to a map[string]any type, returning both the converted map and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// K must be a string type.
//
// Example:
//
//	result, err := StringAnyMapE[string](map[string]interface{}{"key": "value"}) // returns map[string]interface{}{"key": "value"}, nil
//	result, err := StringAnyMapE[string]("invalid") // returns nil, error
func StringAnyMapE[K ~string](o any) (map[K]any, error) {
	return mapE[map[K]any](o, StringE[K], func(o any) (any, error) { return o, nil })
}

// StringStringMap casts an interface to a map[string]string type, ignoring any conversion errors.
// It returns an empty map if conversion fails.
// K and V must be string types.
//
// Example:
//
//	result := StringStringMap[string, string](map[string]string{"key": "value"}) // returns map[string]string{"key": "value"}
func StringStringMap[K ~string, V ~string](o any) map[K]V {
	v, _ := StringStringMapE[K, V](o)
	return v
}

// StringStringMapE casts an interface to a map[string]string type, returning both the converted map and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// K and V must be string types.
//
// Example:
//
//	result, err := StringStringMapE[string, string](map[string]string{"key": "value"}) // returns map[string]string{"key": "value"}, nil
//	result, err := StringStringMapE[string, string]("invalid") // returns nil, error
func StringStringMapE[K ~string, V ~string](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], StringE[V])
}

// StringBoolMap casts an interface to a map[string]bool type, ignoring any conversion errors.
// It returns an empty map if conversion fails.
// K must be a string type and V must be a boolean type.
//
// Example:
//
//	result := StringBoolMap[string, bool](map[string]bool{"key": true}) // returns map[string]bool{"key": true}
func StringBoolMap[K ~string, V ~bool](o any) map[K]V {
	v, _ := StringBoolMapE[K, V](o)
	return v
}

// StringBoolMapE casts an interface to a map[string]bool type, returning both the converted map and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// K must be a string type and V must be a boolean type.
//
// Example:
//
//	result, err := StringBoolMapE[string, bool](map[string]bool{"key": true}) // returns map[string]bool{"key": true}, nil
//	result, err := StringBoolMapE[string, bool]("invalid") // returns nil, error
func StringBoolMapE[K ~string, V ~bool](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], BoolE[V])
}

// StringFloatMap casts an interface to a map[string]float type, ignoring any conversion errors.
// It returns an empty map if conversion fails.
// K must be a string type and V must be a floating-point type.
//
// Example:
//
//	result := StringFloatMap[string, float64](map[string]float64{"key": 3.14}) // returns map[string]float64{"key": 3.14}
func StringFloatMap[K ~string, V constraints.Float](o any) map[K]V {
	v, _ := StringFloatMapE[K, V](o)
	return v
}

// StringFloatMapE casts an interface to a map[string]float type, returning both the converted map and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// K must be a string type and V must be a floating-point type.
//
// Example:
//
//	result, err := StringFloatMapE[string, float64](map[string]float64{"key": 3.14}) // returns map[string]float64{"key": 3.14}, nil
//	result, err := StringFloatMapE[string, float64]("invalid") // returns nil, error
func StringFloatMapE[K ~string, V constraints.Float](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], FloatE[V])
}

// StringIntMap casts an interface to a map[string]int type, ignoring any conversion errors.
// It returns an empty map if conversion fails.
// K must be a string type and V must be a signed integer type.
//
// Example:
//
//	result := StringIntMap[string, int64](map[string]int64{"key": 42}) // returns map[string]int64{"key": 42}
func StringIntMap[K ~string, V constraints.Signed](o any) map[K]V {
	v, _ := StringIntMapE[K, V](o)
	return v
}

// StringIntMapE casts an interface to a map[string]int type, returning both the converted map and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// K must be a string type and V must be a signed integer type.
//
// Example:
//
//	result, err := StringIntMapE[string, int64](map[string]int64{"key": 42}) // returns map[string]int64{"key": 42}, nil
//	result, err := StringIntMapE[string, int64]("invalid") // returns nil, error
func StringIntMapE[K ~string, V constraints.Signed](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], IntE[V])
}

// StringUintMap casts an interface to a map[string]uint type, ignoring any conversion errors.
// It returns an empty map if conversion fails.
// K must be a string type and V must be an unsigned integer type.
//
// Example:
//
//	result := StringUintMap[string, uint64](map[string]uint64{"key": 42}) // returns map[string]uint64{"key": 42}
func StringUintMap[K ~string, V constraints.Unsigned](o any) map[K]V {
	v, _ := StringUintMapE[K, V](o)
	return v
}

// StringUintMapE casts an interface to a map[string]uint type, returning both the converted map and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// K must be a string type and V must be an unsigned integer type.
//
// Example:
//
//	result, err := StringUintMapE[string, uint64](map[string]uint64{"key": 42}) // returns map[string]uint64{"key": 42}, nil
//	result, err := StringUintMapE[string, uint64]("invalid") // returns nil, error
func StringUintMapE[K ~string, V constraints.Unsigned](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], UintE[V])
}

// StringStringSliceMap casts an interface to a map[string][]string type, ignoring any conversion errors.
// It returns an empty map if conversion fails.
//
// Example:
//
//	result := StringStringSliceMap(map[string][]string{"key": {"value1", "value2"}}) // returns map[string][]string{"key": {"value1", "value2"}}
func StringStringSliceMap(o any) map[string][]string {
	v, _ := StringStringSliceMapE(o)
	return v
}

// StringStringSliceMapE casts an interface to a map[string][]string type, returning both the converted map and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
//
// Example:
//
//	result, err := StringStringSliceMapE(map[string][]string{"key": {"value1", "value2"}}) // returns map[string][]string{"key": {"value1", "value2"}}, nil
//	result, err := StringStringSliceMapE("invalid") // returns nil, error
func StringStringSliceMapE(o any) (map[string][]string, error) {
	return mapE[map[string][]string](o, StringE[string], StringSE[[]string])
}

// MapE is a generic function that casts an interface to a map type, returning both the converted map and any error encountered.
// M is the target map type, K is the key type, and V is the value type.
// key is a function that converts keys, and val is a function that converts values.
//
// Example:
//
//	result, err := MapE[map[string]int](map[string]string{"key": "42"}, StringE[string], IntE[int])
//	// returns map[string]int{"key": 42}, nil
func MapE[M ~map[K]V, K comparable, V any](o any, key func(o any) (K, error), val func(o any) (V, error)) (M, error) {
	return mapE[M, K, V](o, key, val)
}

// mapE is the core implementation of map conversion with error handling.
// It uses JSON unmarshaling for string inputs and reflection for map inputs.
// M is the target map type, K is the key type, and V is the value type.
// key is a function that converts keys, and val is a function that converts values.
func mapE[M ~map[K]V, K comparable, V any](o any, key func(o any) (K, error), val func(o any) (V, error)) (M, error) {
	var zero M
	// Handle nil input by returning zero value
	if o == nil {
		return zero, nil
	}

	// Handle string input by JSON unmarshaling
	if s, ok := o.(string); ok {
		res := make(M)
		err := json.Unmarshal([]byte(s), &res)
		if err != nil {
			return failedCastErrValue[M](o, err)
		}
		return res, nil
	}

	// Check if input is a map type
	oType := reflect.TypeOf(o)
	if oType.Kind() != reflect.Map {
		return failedCastValue[M](o)
	}

	// Create result map and populate it by converting each key-value pair
	res := make(M)
	resVal := reflect.ValueOf(res)
	oValue := reflect.ValueOf(o)
	for _, keyVal := range oValue.MapKeys() {
		k, err := key(oValue.MapIndex(keyVal).Interface())
		if err != nil {
			return zero, err
		}
		v, err := val(oValue.MapIndex(keyVal).Interface())
		if err != nil {
			return zero, err
		}
		resVal.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}
	return res, nil
}
