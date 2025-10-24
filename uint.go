package gonv

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Uint converts an interface to a unsigned integer type.
func Uint[N constraints.Unsigned](o any) N {
	v, _ := UintE[N](o)
	return v
}

// UintE converts an interface to a unsigned integer type.
func UintE[E constraints.Unsigned](o any) (E, error) {
	return uintE[E](o)
}

// UintS converts an interface to an unsigned integer slice type.
func UintS[S ~[]E, E constraints.Unsigned](o any) S {
	v, _ := UintSE[S](o)
	return v
}

// UintSE converts an interface to an unsigned integer slice type.
func UintSE[S ~[]E, E constraints.Unsigned](o any) (S, error) {
	return toSliceE[S](o, uintE[E])
}

func uintE[E constraints.Unsigned](o any) (E, error) {
	var zero E
	if o == nil {
		return zero, nil
	}
	// fast path
	switch u := o.(type) {
	case bool:
		if u {
			return 1, nil
		}
		return zero, nil
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
	case string:
		v, err := strconv.ParseUint(trimZeroDecimal(u), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case []byte:
		v, err := strconv.ParseUint(trimZeroDecimal(string(u)), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
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
	case json.Number:
		v, err := u.Int64()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), err
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
	case *durationpb.Duration:
		v := u.AsDuration()
		if v < 0 {
			return failedCastValue[E](o)
		}
		return E(v), nil
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
	default:
		return toUnsignedValueE[E](o)
	}
}

func toUnsignedValueE[E constraints.Unsigned](o any) (E, error) {
	v := indirectValue(reflect.ValueOf(o))
	var zero E
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return zero, nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		u := v.Int()
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return E(v.Uint()), nil
	case reflect.Float64, reflect.Float32:
		u := v.Float()
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	case reflect.String:
		u, err := strconv.ParseUint(trimZeroDecimal(v.String()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		if u < 0 {
			return failedCastValue[E](o)
		}
		return E(u), nil
	default:
		return failedCastValue[E](o)
	}
}
