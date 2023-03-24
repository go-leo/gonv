//go:build jsoniter
// +build jsoniter

package jsonx

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func NewEncoder(w io.Writer) encodingx.Encoder {
	return json.NewEncoder(w)
}

func NewDecoder(r io.Reader) encodingx.Decoder {
	return json.NewDecoder(r)
}
