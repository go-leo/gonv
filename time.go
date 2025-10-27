// Package gonv provides type conversion utilities for Go applications.
// It offers safe and flexible casting between different data types with generic support.
// This file contains functions for converting values to time.Time types.
package gonv

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// DefaultTimeFormat is the default time format used for time string conversions.
// It uses RFC3339 format (YYYY-MM-DDTHH:MM:SSZ).
var DefaultTimeFormat = time.RFC3339

// TimeFormats is a list of time formats supported for parsing time strings.
// It includes common formats like RFC822, RFC850, RFC1123, RFC3339, and several custom formats.
var TimeFormats = []string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	// Handy time stamps.
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,

	"2006-01-02 15:04:05Z07:00",
	"02 Jan 2006",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02T15:04:05",                     // iso8601 without timezone
	"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
	"2006-01-02T15:04:05-0700",                // RFC3339 without timezone hh:mm colon
	"2006-01-02 15:04:05Z0700",                // RFC3339 without T or timezone hh:mm colon
}

// Time casts an interface to a time.Time type, ignoring any conversion errors.
// It returns the zero time value if conversion fails.
// The time is interpreted in UTC location.
//
// Example:
//
//	result := Time("2023-01-01T12:00:00Z") // returns time.Time representing the given timestamp
//	result := Time(1672574400) // returns time.Time representing Unix timestamp
func Time(o any) time.Time {
	v, _ := TimeE(o)
	return v
}

// TimeE casts an interface to a time.Time type, returning both the converted time and any error encountered.
// The time is interpreted in UTC location.
// This function is useful when you need to handle conversion errors explicitly.
//
// Example:
//
//	result, err := TimeE("2023-01-01T12:00:00Z") // returns time.Time and nil
//	result, err := TimeE("invalid") // returns zero time and error
func TimeE(o any) (time.Time, error) {
	return TimeInLocationE(o, time.UTC)
}

// TimeInLocation casts an empty interface to time.Time, interpreting inputs without a timezone
// to be in the given location, or the local timezone if nil.
// It returns the zero time value if conversion fails.
//
// Example:
//
//	loc, _ := time.LoadLocation("America/New_York")
//	result := TimeInLocation("2023-01-01 12:00:00", loc) // returns time.Time in New York timezone
func TimeInLocation(o any, location *time.Location) time.Time {
	v, _ := TimeInLocationE(o, location)
	return v
}

// TimeInLocationE casts an empty interface to time.Time, interpreting inputs without a timezone
// to be in the given location, or the local timezone if nil.
// It returns both the converted time and any error encountered.
// This function is useful when you need to handle conversion errors explicitly.
//
// Example:
//
//	loc, _ := time.LoadLocation("America/New_York")
//	result, err := TimeInLocationE("2023-01-01 12:00:00", loc) // returns time.Time and nil
//	result, err := TimeInLocationE("invalid", loc) // returns zero time and error
func TimeInLocationE(o any, location *time.Location) (time.Time, error) {
	return timeInLocationE(o, location)
}

// timeInLocationE is the core implementation of time conversion with error handling.
// It supports multiple input types and tries to parse them using various time formats.
// The time is interpreted in the given location.
func timeInLocationE(o any, location *time.Location) (time.Time, error) {
	var zero time.Time
	// Handle nil input by returning zero time
	if o == nil {
		return zero, nil
	}

	// Handle different input types
	switch t := o.(type) {
	// String conversion: try parsing with all supported formats
	case string:
		for _, format := range TimeFormats {
			tim, err := time.ParseInLocation(format, t, location)
			if err != nil {
				continue
			}
			return tim, nil
		}
		return failedCastValue[time.Time](o)

	// Byte slice conversion: convert to string and parse
	case []byte:
		ts := string(t)
		for _, format := range TimeFormats {
			tim, err := time.ParseInLocation(format, ts, location)
			if err != nil {
				continue
			}
			return tim, nil
		}
		return failedCastValue[time.Time](o)

	// Native time.Time type: return as is
	case time.Time:
		return t, nil

	// Database driver.Valuer interface support
	case driver.Valuer:
		v, err := t.Value()
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		r, err := timeInLocationE(v, location)
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		return r, nil

	// Protobuf timestamp type support
	case timestamppb.Timestamp:
		return t.AsTime(), nil

	// Protobuf string and bytes wrapper types support
	case *wrapperspb.StringValue:
		r, err := timeInLocationE(t.GetValue(), location)
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		return r, nil
	case *wrapperspb.BytesValue:
		r, err := timeInLocationE(t.GetValue(), location)
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		return r, nil

	// Numeric types: treat as Unix timestamp
	case
		float32, float64,
		int, int64, int32, int16, int8,
		uint, uint64, uint32, uint16, uint8,
		json.Number,
		*durationpb.Duration,
		*wrapperspb.Int64Value, *wrapperspb.Int32Value,
		*wrapperspb.UInt64Value, *wrapperspb.UInt32Value,
		*wrapperspb.DoubleValue, *wrapperspb.FloatValue:
		v, err := IntE[int64](t)
		if err != nil {
			return zero, err
		}
		return time.Unix(v, 0), nil

	// Unsupported types
	default:
		return failedCastValue[time.Time](o)
	}
}
