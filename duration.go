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

// Duration casts an interface to a time.Duration type, ignoring any conversion errors.
// It returns zero duration if conversion fails.
//
// Example:
//   result := Duration("1h30m") // returns 1 hour 30 minutes
//   result := Duration(3600) // returns 3600 nanoseconds
//   result := Duration(int64(3600000000000)) // returns 3600 seconds
func Duration(o any) time.Duration {
	v, _ := DurationE(o)
	return v
}

// DurationE casts an interface to a time.Duration type, returning both the converted value and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
//
// Example:
//   result, err := DurationE("1h30m") // returns 5400000000000, nil
//   result, err := DurationE("invalid") // returns 0, error
func DurationE(o any) (time.Duration, error) {
	return durationE(o)
}

// DurationS casts an interface to a []time.Duration type, ignoring any conversion errors.
// It's designed for converting slice-like data structures to duration slices.
//
// Example:
//   result := DurationS([]string{"1h", "30m", "45s"}) // returns []time.Duration{3600000000000, 1800000000000, 45000000000}
func DurationS(o any) []time.Duration {
	v, _ := DurationSE(o)
	return v
}

// DurationSE casts an interface to a []time.Duration type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors for slice data explicitly.
//
// Example:
//   result, err := DurationSE([]string{"1h", "30m"}) // returns []time.Duration{3600000000000, 1800000000000}, nil
//   result, err := DurationSE([]string{"1h", "invalid"}) // returns nil, error
func DurationSE(o any) ([]time.Duration, error) {
	return toSliceE[[]time.Duration](o, DurationE)
}

// durationE is the core implementation of duration conversion with error handling.
// It uses a fast path approach for common types and falls back to reflection for complex types.
func durationE(o any) (time.Duration, error) {
	// Handle nil input by returning zero duration
	if o == nil {
		var zero time.Duration
		return zero, nil
	}
	
	// Fast path: direct type assertions for common types
	switch d := o.(type) {
	// String conversion using time.ParseDuration
	case string:
		v, err := time.ParseDuration(d)
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return v, nil
		
	// Byte slice conversion by converting to string first
	case []byte:
		v, err := time.ParseDuration(string(d))
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return v, nil
		
	// Stringer interface support for custom types that can be represented as strings
	case fmt.Stringer:
		v, err := time.ParseDuration(d.String())
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return v, nil
		
	// Native time.Duration type
	case time.Duration:
		return d, nil
		
	// Database driver.Valuer interface support
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
		
	// Protobuf duration type support
	case *durationpb.Duration:
		return d.AsDuration(), nil
		
	// Protobuf string wrapper support
	case *wrapperspb.StringValue:
		duration, err := time.ParseDuration(d.GetValue())
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return duration, nil
		
	// Protobuf bytes wrapper support
	case *wrapperspb.BytesValue:
		duration, err := time.ParseDuration(string(d.GetValue()))
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return duration, nil
		
	// Numeric types: convert to int64 first, then create duration
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
		
	// Default case: use reflection-based conversion for complex types
	default:
		// slow path
		return durationVE(o)
	}
}

// durationVE is the reflection-based (slow path) implementation for duration conversion.
// It's used when fast path type assertions fail and more complex type analysis is needed.
func durationVE(o any) (time.Duration, error) {
	// Get the underlying value, dereferencing pointers if necessary
	v := indirectValue(reflect.ValueOf(o))
	
	// Handle different reflection kinds
	switch v.Kind() {
	// Integer types: directly convert to duration (interpreted as nanoseconds)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return time.Duration(v.Int()), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return time.Duration(v.Uint()), nil
		
	// Floating point types: convert to duration (interpreted as nanoseconds)
	case reflect.Float64, reflect.Float32:
		return time.Duration(v.Float()), nil
		
	// String conversion using time.ParseDuration
	case reflect.String:
		dur, err := time.ParseDuration(v.String())
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return dur, nil
		
	// Byte slice conversion (must be []byte)
	case reflect.Slice:
		// Ensure it's a byte slice
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[time.Duration](o)
		}
		dur, err := time.ParseDuration(string(v.Bytes()))
		if err != nil {
			return failedCastErrValue[time.Duration](o, err)
		}
		return dur, nil
		
	// Unsupported types
	default:
		return failedCastValue[time.Duration](o)
	}
}