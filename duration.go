package gonv

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Duration casts an interface to a time.Duration type.
func Duration(o any) time.Duration {
	v, _ := DurationE(o)
	return v
}

// DurationE casts an interface to a time.Duration type.
func DurationE(o any) (time.Duration, error) {
	return durationE(o)
}

// DurationS casts an interface to a []time.Duration type.
func DurationS(o any) []time.Duration {
	v, _ := DurationSE(o)
	return v
}

// DurationSE casts an interface to a []time.Duration type.
func DurationSE(o any) ([]time.Duration, error) {
	return toSliceE[[]time.Duration](o, DurationE)
}

func durationE(o any) (time.Duration, error) {
	if o == nil {
		var zero time.Duration
		return zero, nil
	}
	// fast path
	switch d := o.(type) {
	case string:
		v, err := time.ParseDuration(d)
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return v, nil
	case []byte:
		v, err := time.ParseDuration(string(d))
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return v, nil
	case fmt.Stringer:
		v, err := time.ParseDuration(d.String())
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return v, nil
	case time.Duration:
		return d, nil
	case driver.Valuer:
		v, err := d.Value()
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		r, err := durationE(v)
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return r, nil
	case *durationpb.Duration:
		return d.AsDuration(), nil
	case *wrapperspb.StringValue:
		duration, err := time.ParseDuration(d.GetValue())
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return duration, nil
	case *wrapperspb.BytesValue:
		duration, err := time.ParseDuration(string(d.GetValue()))
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return duration, nil
	case
		float32, float64,
		int, int64, int32, int16, int8,
		uint, uint64, uint32, uint16, uint8,
		json.Number,
		*wrapperspb.DoubleValue,
		*wrapperspb.FloatValue,
		*wrapperspb.Int64Value,
		*wrapperspb.Int32Value,
		*wrapperspb.UInt64Value,
		*wrapperspb.UInt32Value:
		duration, err := intE[time.Duration](o)
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return time.Duration(duration), nil
	default:
		// slow path
		return durationVE(o)
	}
}

func durationVE(o any) (time.Duration, error) {
	v := indirectValue(reflect.ValueOf(o))
	switch v.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return time.Duration(v.Int()), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return time.Duration(v.Uint()), nil
	case reflect.Float64, reflect.Float32:
		return time.Duration(v.Float()), nil
	case reflect.String:
		dur, err := time.ParseDuration(v.String())
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return dur, nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[time.Duration](o)
		}
		dur, err := time.ParseDuration(string(v.Bytes()))
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return dur, nil
	default:
		return failedCastValue[time.Duration](o)
	}
}
