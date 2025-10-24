package gonv

import (
	"database/sql/driver"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// ToDuration casts an interface to a time.Duration type.
func ToDuration(o any) time.Duration {
	v, _ := ToDurationE(o)
	return v
}

// ToDurationE casts an interface to a time.Duration type.
func ToDurationE(o any) (time.Duration, error) {
	return toDurationE(o)
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(o any) []time.Duration {
	v, _ := ToDurationSliceE(o)
	return v
}

// ToDurationSliceE casts an interface to a []time.Duration type.
func ToDurationSliceE(o any) ([]time.Duration, error) {
	return toSliceE[[]time.Duration](o, ToDurationE)
}

func toDurationE(o any) (time.Duration, error) {
	var zero time.Duration
	if o == nil {
		return zero, nil
	}
	// fast path
	switch d := o.(type) {
	case time.Duration:
		return d, nil
	case *durationpb.Duration:
		return d.AsDuration(), nil
	case string:
		duration, err := time.ParseDuration(d)
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return duration, nil
	case []byte:
		duration, err := time.ParseDuration(string(d))
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return duration, nil
	case int, int64, int32, int16, int8,
		uint, uint64, uint32, uint16, uint8,
		float32, float64,
		int64er, float64er,
		driver.Valuer,
		*wrapperspb.DoubleValue,
		*wrapperspb.FloatValue,
		*wrapperspb.Int64Value,
		*wrapperspb.UInt64Value,
		*wrapperspb.Int32Value,
		*wrapperspb.UInt32Value,
		*wrapperspb.StringValue,
		*wrapperspb.BytesValue:
		duration, err := intE[time.Duration](o)
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return time.Duration(duration), nil
	default:
		// slow path
		return toDurationValueE(o)
	}
}

func toDurationValueE(o any) (time.Duration, error) {
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
	default:
		return failedCastValue[time.Duration](o)
	}
}
