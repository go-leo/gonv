package gonv

import (
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// String casts an interface to a string type, ignoring any conversion errors.
// It returns an empty string if conversion fails.
// E must be a string type.
//
// Example:
//
//	result := String[string](42) // returns "42"
//	result := String[string](true) // returns "true"
//	result := String[string](3.14) // returns "3.14"
func String[E ~string](o any) E {
	v, _ := StringE[E](o)
	return v
}

// StringE casts an interface to a string type, returning both the converted string and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// E must be a string type.
//
// Example:
//
//	result, err := StringE[string](42) // returns "42", nil
//	result, err := StringE[string](nil) // returns "", nil
func StringE[E ~string](o any) (E, error) {
	return stringE[E](o)
}

// StringS casts an interface to a []string type, ignoring any conversion errors.
// It returns an empty slice if conversion fails.
// S is a slice type with elements of string type.
// E must be a string type.
//
// Example:
//
//	result := StringS[[]string, string]([]int{1, 2, 3}) // returns []string{"1", "2", "3"}
func StringS[S ~[]E, E ~string](o any) S {
	v, _ := StringSE[S](o)
	return v
}

// StringSE casts an interface to a []string type, returning both the converted slice and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
// S is a slice type with elements of string type.
// E must be a string type.
//
// Example:
//
//	result, err := StringSE[[]string, string]([]int{1, 2, 3}) // returns []string{"1", "2", "3"}, nil
//	result, err := StringSE[[]string, string]("not a slice") // returns nil, error
func StringSE[S ~[]E, E ~string](o any) (S, error) {
	return toSliceE[S](o, stringE[E])
}

// stringE is the core implementation of string conversion with error handling.
// It uses a fast path approach for common types and falls back to reflection for complex types.
// E must be a string type.
func stringE[E ~string](o any) (E, error) {
	var zero E
	// Handle nil input by returning zero value
	if o == nil {
		return zero, nil
	}

	// Fast path: direct type assertions for common types
	switch s := o.(type) {
	// Boolean conversion using strconv.FormatBool
	case bool:
		return E(strconv.FormatBool(s)), nil

	// Floating-point conversion using strconv.FormatFloat
	case float64:
		return E(strconv.FormatFloat(s, 'f', -1, 64)), nil
	case float32:
		return E(strconv.FormatFloat(float64(s), 'f', -1, 32)), nil

	// Integer conversion using strconv functions
	case int:
		return E(strconv.Itoa(s)), nil
	case int64:
		return E(strconv.FormatInt(s, 10)), nil
	case int32:
		return E(strconv.FormatInt(int64(s), 10)), nil
	case int16:
		return E(strconv.FormatInt(int64(s), 10)), nil
	case int8:
		return E(strconv.FormatInt(int64(s), 10)), nil
	case uint:
		return E(strconv.FormatUint(uint64(s), 10)), nil
	case uint64:
		return E(strconv.FormatUint(s, 10)), nil
	case uint32:
		return E(strconv.FormatUint(uint64(s), 10)), nil
	case uint16:
		return E(strconv.FormatUint(uint64(s), 10)), nil
	case uint8:
		return E(strconv.FormatUint(uint64(s), 10)), nil

	// String types: direct conversion
	case string:
		return E(s), nil
	case []byte:
		return E(string(s)), nil

	// HTML template types: convert to string
	case template.HTML:
		return E(string(s)), nil
	case template.URL:
		return E(string(s)), nil
	case template.JS:
		return E(string(s)), nil
	case template.CSS:
		return E(string(s)), nil
	case template.HTMLAttr:
		return E(string(s)), nil

	// JSON number: use String() method
	case json.Number:
		return E(s.String()), nil

	// Time types: use String() method or Format() with DefaultTimeFormat
	case time.Weekday:
		return E(s.String()), nil
	case time.Month:
		return E(s.String()), nil
	case time.Duration:
		return E(s.String()), nil
	case time.Time:
		return E(s.Format(DefaultTimeFormat)), nil

	// Protobuf types support
	case *durationpb.Duration:
		return E(s.AsDuration().String()), nil
	case *timestamppb.Timestamp:
		return E(s.AsTime().Format(DefaultTimeFormat)), nil

	// Protobuf wrapper types support
	case *wrapperspb.BoolValue:
		return E(strconv.FormatBool(s.GetValue())), nil
	case *wrapperspb.DoubleValue:
		return E(strconv.FormatFloat(s.GetValue(), 'f', -1, 64)), nil
	case *wrapperspb.FloatValue:
		return E(strconv.FormatFloat(float64(s.GetValue()), 'f', -1, 32)), nil
	case *wrapperspb.Int64Value:
		return E(strconv.FormatInt(s.GetValue(), 10)), nil
	case *wrapperspb.Int32Value:
		return E(strconv.FormatInt(int64(s.GetValue()), 10)), nil
	case *wrapperspb.UInt64Value:
		return E(strconv.FormatUint(s.GetValue(), 10)), nil
	case *wrapperspb.UInt32Value:
		return E(strconv.FormatUint(uint64(s.GetValue()), 10)), nil
	case *wrapperspb.StringValue:
		return E(s.GetValue()), nil
	case *wrapperspb.BytesValue:
		return E(s.GetValue()), nil

	// Database driver.Valuer interface support
	case driver.Valuer:
		v, err := s.Value()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		r, err := stringE[E](v)
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return r, nil

	// Text and JSON marshaling interfaces
	case encoding.TextMarshaler:
		v, err := s.MarshalText()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(string(v)), nil

	case json.Marshaler:
		v, err := s.MarshalJSON()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return E(string(v)), nil

	// Stringer and error interfaces
	case fmt.Stringer:
		return E(s.String()), nil
	case error:
		return E(s.Error()), nil

	// Default case: use reflection-based conversion for complex types
	default:
		// slow path
		return stringVE[E](o)
	}
}

// stringVE is the reflection-based (slow path) implementation for string conversion.
// It's used when fast path type assertions fail and more complex type analysis is needed.
// E must be a string type.
func stringVE[E ~string](o any) (E, error) {
	// Get the underlying value, dereferencing pointers if necessary
	v := indirectValue(reflect.ValueOf(o))

	// Handle different reflection kinds
	switch v.Kind() {
	// Boolean conversion using strconv.FormatBool
	case reflect.Bool:
		return E(strconv.FormatBool(v.Bool())), nil

	// Integer conversion using strconv.FormatInt/FormatUint
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return E(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return E(strconv.FormatUint(v.Uint(), 10)), nil

	// Floating-point conversion using strconv.FormatFloat
	case reflect.Float64, reflect.Float32:
		return E(strconv.FormatFloat(v.Float(), 'f', -1, 64)), nil

	// String conversion
	case reflect.String:
		return E(v.String()), nil

	// Byte slice conversion
	case reflect.Slice:
		// Ensure it's a byte slice
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return failedCastValue[E](o)
		}
		return E(string(v.Bytes())), nil

	// Unsupported types
	default:
		return failedCastValue[E](o)
	}
}
