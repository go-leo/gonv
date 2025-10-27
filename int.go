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
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Int converts an interface to a signed integer type.
func Int[E constraints.Signed](o any) E {
	v, _ := IntE[E](o)
	return v
}

// IntE converts an interface to a signed integer type.
func IntE[E constraints.Signed](o any) (E, error) {
	return intE[E](o)
}

// IntS converts an interface to a signed integer slice type.
func IntS[S ~[]E, E constraints.Signed](o any) S {
	v, _ := IntSE[S](o)
	return v
}

// IntSE converts an interface to a signed integer slice type.
func IntSE[S ~[]E, E constraints.Signed](o any) (S, error) {
	return toSliceE[S](o, IntE[E])
}

func intE[E constraints.Signed](o any) (E, error) {
	var zero E
	if o == nil {
		return zero, nil
	}
	switch s := o.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return zero, nil
	case float64:
		return E(s), nil
	case float32:
		return E(s), nil
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
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case []byte:
		v, err := strconv.ParseInt(trimZeroDecimal(string(s)), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case fmt.Stringer:
		v, err := strconv.ParseInt(trimZeroDecimal(s.String()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case json.Number:
		v, err := s.Int64()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case time.Weekday:
		return E(s), nil
	case time.Month:
		return E(s), nil
	case time.Duration:
		return E(s), nil
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
	case *durationpb.Duration:
		return E(s.AsDuration()), nil
	case *timestamppb.Timestamp:
		return E(s.AsTime().UnixMilli()), nil
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
	default:
		// slow path
		return toSignedValueE[E](o)
	}
}

func toSignedValueE[E constraints.Signed](o any) (E, error) {
	var zero E
	v := indirectValue(reflect.ValueOf(o))
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return zero, nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return E(v.Int()), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return E(v.Uint()), nil
	case reflect.Float64, reflect.Float32:
		return E(v.Float()), nil
	case reflect.String:
		i, err := strconv.ParseInt(trimZeroDecimal(v.String()), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(i), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		i, err := strconv.ParseInt(trimZeroDecimal(string(v.Bytes())), 0, 0)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(i), nil
	default:
		return failedCastValue[E](o)
	}
}
