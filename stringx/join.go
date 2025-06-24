package stringx

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func JoinInt[E constraints.Signed](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatInt(int64(e), 10))
	}
	return strings.Join(s, sep)
}

func JoinUint[E constraints.Unsigned](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatUint(uint64(e), 10))
	}
	return strings.Join(s, sep)
}

func JoinFloat[E constraints.Float](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatFloat(float64(e), 'f', -1, 64))
	}
	return strings.Join(s, sep)
}

func JoinBool[E ~bool](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatBool(bool(e)))
	}
	return strings.Join(s, sep)
}
