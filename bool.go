package gonv

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Bool casts an interface to a bool type.
func Bool[E ~bool](o any) E {
	v, _ := BoolE[E](o)
	return v
}

// BoolE casts an interface to a bool type.
func BoolE[E ~bool](o any) (E, error) {
	return boolE[E](o)
}

// BoolS casts an interface to a []bool type.
func BoolS[S ~[]E, E ~bool](o any) S {
	v, _ := BoolSE[S](o)
	return v
}

// BoolSE casts an interface to a []bool type.
func BoolSE[S ~[]E, E ~bool](o any) (S, error) {
	return toSliceE[S](o, boolE[E])
}

func boolE[E ~bool](o any) (E, error) {
	if o == nil {
		var zero E
		return zero, nil
	}
	// fast path
	switch b := o.(type) {
	case bool:
		return E(b), nil
	case string:
		v, err := strconv.ParseBool(b)
		if err != nil {
			return failedCastErrValue[E](b, err)
		}
		return E(v), err
	case []byte:
		v, err := strconv.ParseBool(string(b))
		if err != nil {
			return failedCastErrValue[E](b, err)
		}
		return E(v), err
	case fmt.Stringer:
		v, err := strconv.ParseBool(b.String())
		if err != nil {
			return failedCastErrValue[E](b, err)
		}
		return E(v), err
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
		return n != 0, nil
	default:
		// slow path
		return boolVE[E](o)
	}
}

func boolVE[E ~bool](o any) (E, error) {
	v := indirectValue(reflect.ValueOf(o))
	switch v.Kind() {
	case reflect.Bool:
		return E(v.Bool()), nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return v.Int() != 0, nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return v.Uint() != 0, nil
	case reflect.Float64, reflect.Float32:
		return v.Float() != 0, nil
	case reflect.String:
		b, err := strconv.ParseBool(v.String())
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(b), err
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		b, err := strconv.ParseBool(string(v.Bytes()))
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(b), err
	default:
		return failedCastValue[E](o)
	}
}
