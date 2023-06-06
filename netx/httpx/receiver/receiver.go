package receiver

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-leo/gox/convx"
	"github.com/go-leo/gox/encodingx/jsonx"
	"github.com/go-leo/gox/encodingx/xmlx"
	"github.com/go-leo/gox/errorx"
	"github.com/go-leo/gox/slicex"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

type Receiver interface {
	Response() *http.Response
	Status() string
	StatusCode() int
	Proto() string
	ProtoMajor() int
	ProtoMinor() int
	ContentLength() int64
	TransferEncoding() []string
	Headers() http.Header
	Trailers() http.Header
	Cookies() []*http.Cookie
	BytesBody() ([]byte, error)
	TextBody() (string, error)
	ObjectBody(body any, unmarshal func([]byte, any) error) error
	JSONBody(body any) error
	XMLBody(body any) error
	ProtobufBody(body proto.Message) error
	GobBody(body any) error
	WriterBody(file io.Writer) error
}

type receiver struct {
	resp      *http.Response
	err       error
	bodyBytes []byte
	read      bool
}

func (r *receiver) Response() *http.Response {
	return r.resp
}

func (r *receiver) Status() string {
	return r.resp.Status
}

func (r *receiver) StatusCode() int {
	return r.resp.StatusCode
}

func (r *receiver) Proto() string {
	return r.resp.Proto
}

func (r *receiver) ProtoMajor() int {
	return r.resp.ProtoMajor
}

func (r *receiver) ProtoMinor() int {
	return r.resp.ProtoMinor
}

func (r *receiver) ContentLength() int64 {
	return r.resp.ContentLength
}

func (r *receiver) TransferEncoding() []string {
	return r.resp.TransferEncoding
}

func (r *receiver) Headers() http.Header {
	return r.resp.Header
}

func (r *receiver) Trailers() http.Header {
	return r.resp.Trailer
}

func (r *receiver) Cookies() []*http.Cookie {
	return r.resp.Cookies()
}

func (r *receiver) BytesBody() ([]byte, error) {
	if r.read {
		return r.bodyBytes, nil
	}
	r.read = true
	body, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return nil, err
	}
	defer errorx.Silence(r.resp.Body.Close())
	r.bodyBytes = body
	return body, nil
}

func (r *receiver) TextBody() (string, error) {
	bodyBytes, err := r.BytesBody()
	if err != nil {
		return "", err
	}
	return convx.BytesToString(bodyBytes), nil
}

func (r *receiver) ObjectBody(body any, unmarshal func([]byte, any) error) error {
	bodyBytes, err := r.BytesBody()
	if err != nil {
		return err
	}
	if err := unmarshal(bodyBytes, body); err != nil {
		err = fmt.Errorf("failed to unmarshal body, body is %s, %w", r.bodyBytes, err)
	}
	return nil
}

func (r *receiver) JSONBody(body any) error {
	return r.ObjectBody(body, jsonx.Unmarshal)
}

func (r *receiver) XMLBody(body any) error {
	return r.ObjectBody(body, xmlx.Unmarshal)
}

func (r *receiver) ProtobufBody(body proto.Message) error {
	unmarshal := func(data []byte, v any) error { return proto.Unmarshal(data, v.(proto.Message)) }
	return r.ObjectBody(body, unmarshal)
}

func (r *receiver) GobBody(body any) error {
	unmarshal := func(data []byte, v any) error { return gob.NewDecoder(bytes.NewReader(data)).Decode(v) }
	return r.ObjectBody(body, unmarshal)
}

func (r *receiver) WriterBody(file io.Writer) error {
	r.read = true
	_, err := io.Copy(file, r.resp.Body)
	return err
}

func NewReceiver(resp *http.Response, errs ...error) Receiver {
	var err error
	if slicex.IsNotEmpty(errs) {
		err = errs[0]
	}
	r := &receiver{
		err:  err,
		resp: resp,
	}
	return r
}
