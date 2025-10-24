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
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// String casts an interface to a string type.
func String[E ~string](o any) E {
	v, _ := StringE[E](o)
	return v
}

// StringE casts an interface to a string type.
func StringE[E ~string](o any) (E, error) {
	return stringE[E](o)
}

// StringS casts an interface to a []string type.
func StringS[S ~[]E, E ~string](o any) S {
	v, _ := StringSE[S](o)
	return v
}

// StringSE casts an interface to a []string type.
func StringSE[S ~[]E, E ~string](o any) (S, error) {
	return toSliceE[S](o, stringE[E])
}

func stringE[E ~string](o any) (E, error) {
	var zero E
	if o == nil {
		return zero, nil
	}
	// fast path
	switch s := o.(type) {
	case bool:
		return E(strconv.FormatBool(s)), nil
	case float64:
		return E(strconv.FormatFloat(s, 'f', -1, 64)), nil
	case float32:
		return E(strconv.FormatFloat(float64(s), 'f', -1, 32)), nil
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
	case string:
		return E(s), nil
	case []byte:
		return E(string(s)), nil
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
	case fmt.Stringer:
		return E(s.String()), nil
	case error:
		return E(s.Error()), nil
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
	case json.Number:
		return E(s.String()), nil
	case driver.Valuer:
		v, err := s.Value()
		if err != nil {
			return failedCastErrValue[E](o, err)
		}
		return stringE[E](v)
	case time.Duration:
		return E(s.String()), nil
	case *durationpb.Duration:
		return E(s.AsDuration().String()), nil
	case *wrapperspb.BoolValue:
		return E(strconv.FormatBool(s.GetValue())), nil
	case *wrapperspb.Int64Value:
		return E(strconv.FormatInt(s.GetValue(), 10)), nil
	case *wrapperspb.Int32Value:
		return E(strconv.FormatInt(int64(s.GetValue()), 10)), nil
	case *wrapperspb.UInt64Value:
		return E(strconv.FormatUint(s.GetValue(), 10)), nil
	case *wrapperspb.UInt32Value:
		return E(strconv.FormatUint(uint64(s.GetValue()), 10)), nil
	case *wrapperspb.DoubleValue:
		return E(strconv.FormatFloat(s.GetValue(), 'f', -1, 64)), nil
	case *wrapperspb.FloatValue:
		return E(strconv.FormatFloat(float64(s.GetValue()), 'f', -1, 32)), nil
	case *wrapperspb.StringValue:
		return E(s.GetValue()), nil
	case *wrapperspb.BytesValue:
		return E(s.GetValue()), nil
	default:
		// slow path
		return stringValueE[E](o)
	}
}

func stringValueE[E ~string](o any) (E, error) {
	v := indirectValue(reflect.ValueOf(o))
	switch v.Kind() {
	case reflect.Bool:
		return E(strconv.FormatBool(v.Bool())), nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return E(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return E(strconv.FormatUint(v.Uint(), 10)), nil
	case reflect.Float64, reflect.Float32:
		return E(strconv.FormatFloat(v.Float(), 'f', -1, 64)), nil
	case reflect.String:
		return E(v.String()), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return E(string(v.Bytes())), nil
		}
		return failedCastValue[E](o)
	default:
		return failedCastValue[E](o)
	}
}
