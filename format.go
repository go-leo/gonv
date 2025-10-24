package gonv

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/constraints"
)

// FormatBool takes a boolean-type generic parameter `b`, converts it to a string, and returns
// the string.
func FormatBool[Bool ~bool](b Bool) string {
	return strconv.FormatBool(bool(b))
}

// FormatUint converts an unsigned integer to a string representation in a specified base.
// It does this by first converting the integer to a uint64 and then using the strconv.FormatUint
// function to format it as a string in the desired base.
func FormatUint[Unsigned constraints.Unsigned](i Unsigned, base int) string {
	return strconv.FormatUint(uint64(i), base)
}

func FormatInt[Signed constraints.Signed](i Signed, base int) string {
	return strconv.FormatInt(int64(i), base)
}

func FormatFloat[Float constraints.Float](f Float, fmt byte, prec, bitSize int) string {
	return strconv.FormatFloat(float64(f), fmt, prec, bitSize)
}

func FormatBoolSlice[Bool ~bool](s []Bool) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, b := range s {
		r = append(r, FormatBool(b))
	}
	return r
}

func FormatUintSlice[Unsigned constraints.Unsigned](s []Unsigned, base int) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, i := range s {
		r = append(r, FormatUint(i, base))
	}
	return r
}

func FormatIntSlice[Signed constraints.Signed](s []Signed, base int) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, i := range s {
		r = append(r, FormatInt(i, base))
	}
	return r
}

func FormatFloatSlice[Float constraints.Float](s []Float, fmt byte, prec, bitSize int) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, f := range s {
		r = append(r, FormatFloat(float64(f), fmt, prec, bitSize))
	}
	return r
}

func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func ParseInt[Signed constraints.Signed](s string, base int, bitSize int) (Signed, error) {
	i, err := strconv.ParseInt(s, base, bitSize)
	return Signed(i), err
}

func ParseUint[Unsigned constraints.Unsigned](s string, base int, bitSize int) (Unsigned, error) {
	i, err := strconv.ParseUint(s, base, bitSize)
	return Unsigned(i), err
}

func ParseFloat[Float constraints.Float](s string, bitSize int) (Float, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	return Float(f), err
}

func ParseBoolSlice(s []string) ([]bool, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]bool, 0, len(s))
	for _, str := range s {
		b, err := strconv.ParseBool(str)
		if err != nil {
			return nil, err
		}
		r = append(r, b)
	}
	return r, nil
}

func ParseIntSlice[Signed constraints.Signed](s []string, base int, bitSize int) ([]Signed, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]Signed, 0, len(s))
	for _, str := range s {
		i, err := ParseInt[Signed](str, base, bitSize)
		if err != nil {
			return nil, err
		}
		r = append(r, i)
	}
	return r, nil
}

func ParseUintSlice[Unsigned constraints.Unsigned](s []string, base int, bitSize int) ([]Unsigned, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]Unsigned, 0, len(s))
	for _, str := range s {
		i, err := ParseUint[Unsigned](str, base, bitSize)
		if err != nil {
			return nil, err
		}
		r = append(r, i)
	}
	return r, nil
}

func ParseFloatSlice[Float constraints.Float](s []string, bitSize int) ([]Float, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]Float, 0, len(s))
	for _, str := range s {
		f, err := ParseFloat[Float](str, bitSize)
		if err != nil {
			return nil, err
		}
		r = append(r, f)
	}
	return r, nil
}

func ParseBytesSlice(s []string) [][]byte {
	if s == nil {
		return nil
	}
	r := make([][]byte, 0, len(s))
	for _, str := range s {
		r = append(r, []byte(str))
	}
	return r
}

var quotePool = sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, 16)) }}

// Quote quotes the string.
func Quote[E ~string](e E, quote string) E {
	buffer := quotePool.Get().(*bytes.Buffer)
	defer quotePool.Put(buffer)
	buffer.Reset()
	buffer.WriteString(quote)
	buffer.WriteString(string(e))
	buffer.WriteString(quote)
	return E(buffer.String())
}

func quoteV2[E ~string](e E, quote string) E {
	buffer := quotePool.Get().(*bytes.Buffer)
	defer quotePool.Put(buffer)
	buffer.Reset()
	_, _ = buffer.WriteString(fmt.Sprintf("%s%s%s", quote, e, quote))
	return E(buffer.String())
}

func quoteV3[E ~string](e E, quote string) E {
	return E(strings.Join([]string{quote, string(e), quote}, ""))
}

// QuoteSlice quotes each string in the slice.
func QuoteSlice[S ~[]E, E ~string](s S, quote string) S {
	if s == nil {
		return s
	}
	r := make(S, 0, len(s))
	for _, e := range s {
		r = append(r, Quote(e, quote))
	}
	return r
}
