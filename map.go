package gonv

import (
	"encoding/json"
	"reflect"

	"golang.org/x/exp/constraints"
)

// StringAnyMap casts an interface to a map[string]any type.
func StringAnyMap[K ~string](o any) map[K]any {
	v, _ := StringAnyMapE[K](o)
	return v
}

// StringAnyMapE casts an interface to a map[string]any type.
func StringAnyMapE[K ~string](o any) (map[K]any, error) {
	return mapE[map[K]any](o, StringE[K], func(o any) (any, error) { return o, nil })
}

// StringStringMap casts an interface to a map[string]string type.
func StringStringMap[K ~string, V ~string](o any) map[K]V {
	v, _ := StringStringMapE[K, V](o)
	return v
}

// StringStringMapE casts an interface to a map[string]string type.
func StringStringMapE[K ~string, V ~string](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], StringE[V])
}

// StringBoolMap casts an interface to a map[string]bool type.
func StringBoolMap[K ~string, V ~bool](o any) map[K]V {
	v, _ := StringBoolMapE[K, V](o)
	return v
}

// StringBoolMapE casts an interface to a map[string]bool type.
func StringBoolMapE[K ~string, V ~bool](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], BoolE[V])
}

// StringIntMap casts an interface to a map[string]int type.
func StringFloatMap[K ~string, V constraints.Float](o any) map[K]V {
	v, _ := StringFloatMapE[K, V](o)
	return v
}

// StringIntMapE casts an interface to a map[string]int{} type.
func StringFloatMapE[K ~string, V constraints.Float](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], FloatE[V])
}

// StringIntMap casts an interface to a map[string]int type.
func StringIntMap[K ~string, V constraints.Signed](o any) map[K]V {
	v, _ := StringIntMapE[K, V](o)
	return v
}

// StringIntMapE casts an interface to a map[string]int{} type.
func StringIntMapE[K ~string, V constraints.Signed](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], IntE[V])
}

// StringIntMap casts an interface to a map[string]int type.
func StringUintMap[K ~string, V constraints.Unsigned](o any) map[K]V {
	v, _ := StringUintMapE[K, V](o)
	return v
}

// StringIntMapE casts an interface to a map[string]int{} type.
func StringUintMapE[K ~string, V constraints.Unsigned](o any) (map[K]V, error) {
	return mapE[map[K]V](o, StringE[K], UintE[V])
}

// StringStringSliceMap casts an interface to a map[string][]string type.
func StringStringSliceMap(o any) map[string][]string {
	v, _ := StringStringSliceMapE(o)
	return v
}

// StringStringSliceMapE casts an interface to a map[string][]string type.
func StringStringSliceMapE(o any) (map[string][]string, error) {
	return mapE[map[string][]string](o, StringE[string], StringSE[[]string])
}

func MapE[M ~map[K]V, K comparable, V any](o any, key func(o any) (K, error), val func(o any) (V, error)) (M, error) {
	return mapE[M, K, V](o, key, val)
}

func mapE[M ~map[K]V, K comparable, V any](o any, key func(o any) (K, error), val func(o any) (V, error)) (M, error) {
	var zero M
	if o == nil {
		return zero, nil
	}
	if s, ok := o.(string); ok {
		res := make(M)
		err := json.Unmarshal([]byte(s), &res)
		if err != nil {
			return failedCastErrValue[M](o, err)
		}
		return res, nil
	}
	oType := reflect.TypeOf(o)
	if oType.Kind() != reflect.Map {
		return failedCastValue[M](o)
	}

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
