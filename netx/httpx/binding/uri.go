package binding

import (
	"github.com/go-leo/gox/encodingx/formx"
)

func Uri(m map[string][]string, obj any, tag string) error {
	return formx.Unmarshal(m, obj, tag)
}
