package binding

import (
	"errors"
	"github.com/go-leo/gox/encodingx/jsonx"
	"net/http"
)

func JSON(req *http.Request, obj any, useNumber bool, disallowUnknownFields bool) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	decoder := jsonx.NewDecoder(req.Body)
	if useNumber {
		u, ok := decoder.(interface{ UseNumber() })
		if ok {
			u.UseNumber()
		}
	}
	if disallowUnknownFields {
		d, ok := decoder.(interface{ DisallowUnknownFields() })
		if ok {
			d.DisallowUnknownFields()
		}
	}
	return decoder.Decode(obj)
}
