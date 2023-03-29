package render

import (
	"net/http"

	"github.com/go-leo/gox/encodingx/protobufx"
)

// ProtoBuf marshals the given interface object and writes data with custom ContentType.
func ProtoBuf(w http.ResponseWriter, data any) error {
	writeContentType(w, []string{"application/x-protobuf"})
	bytes, err := protobufx.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}
