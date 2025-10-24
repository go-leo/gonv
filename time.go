package gonv

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

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

// Time casts an interface to a time.Time type.
func Time(o any) time.Time {
	v, _ := TimeE(o)
	return v
}

// TimeE casts an interface to a time.Time type.
func TimeE(o any) (time.Time, error) {
	return TimeInLocationE(o, time.UTC)
}

// TimeInLocation casts an empty interface to time.Time,
func TimeInLocation(o any, location *time.Location) time.Time {
	v, _ := TimeInLocationE(o, location)
	return v
}

// TimeInLocationE casts an empty interface to time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func TimeInLocationE(o any, location *time.Location) (time.Time, error) {
	return toTimeInLocationE(o, location)
}

func toTimeInLocationE(o any, location *time.Location) (time.Time, error) {
	var zero time.Time
	if o == nil {
		return zero, nil
	}
	switch t := o.(type) {
	case time.Time:
		return t, nil
	case timestamppb.Timestamp:
		return t.AsTime(), nil
	case string:
		for _, format := range TimeFormats {
			tim, err := time.ParseInLocation(format, t, location)
			if err != nil {
				continue
			}
			return tim, nil
		}
		return failedCastValue[time.Time](o)
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
	case *wrapperspb.StringValue:
		r, err := toTimeInLocationE(t.GetValue(), location)
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		return r, nil
	case *wrapperspb.BytesValue:
		r, err := toTimeInLocationE(t.GetValue(), location)
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		return r, nil
	case driver.Valuer:
		v, err := t.Value()
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		r, err := toTimeInLocationE(v, location)
		if err != nil {
			return failedCastErrValue[time.Time](o, err)
		}
		return r, nil
	case int, int64, int32, int16, int8,
		uint, uint64, uint32, uint16, uint8,
		float32, float64,
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
	default:
		return failedCastValue[time.Time](o)
	}
}
