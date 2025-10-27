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
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Float converts an interface to a floating-point type.
func Float[E constraints.Float](o any) E {
	v, _ := FloatE[E](o)
	return v
}

// FloatE converts an interface to a floating-point type.
func FloatE[E constraints.Float](o any) (E, error) {
	return floatE[E](o)
}

// FloatS converts an interface to a floating-point slice type.
func FloatS[S ~[]E, E constraints.Float](o any) S {
	v, _ := FloatSE[S](o)
	return v
}

// FloatSE converts an interface to a floating-point slice type.
func FloatSE[S ~[]E, E constraints.Float](o any) (S, error) {
	return toSliceE[S](o, floatE[E])
}

func floatE[E constraints.Float](o any) (E, error) {
	var zero E
	if o == nil {
		return zero, nil
	}
	// fast path
	switch f := o.(type) {
	case bool:
		if f {
			return 1, nil
		}
		return zero, nil
	case float64:
		return E(f), nil
	case float32:
		return E(f), nil
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
	case string:
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case []byte:
		v, err := strconv.ParseFloat(string(f), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case fmt.Stringer:
		v, err := strconv.ParseFloat(f.String(), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case json.Number:
		v, err := f.Float64()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(v), nil
	case time.Weekday:
		return E(f), nil
	case time.Month:
		return E(f), nil
	case time.Duration:
		return E(f), nil
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
	case *durationpb.Duration:
		return E(f.AsDuration()), nil
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
	default:
		// slow path
		return floatVE[E](o)
	}
}

func floatVE[E constraints.Float](o any) (E, error) {
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
		f, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(f), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		f, err := strconv.ParseFloat(string(v.Bytes()), 64)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(f), nil
	default:
		return failedCastValue[E](o)
	}
}
