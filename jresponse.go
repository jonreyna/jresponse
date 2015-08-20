package jresponse

import (
	"io"
)

type JResponseWriter interface {
	WriteXMLTo(w io.Writer) (n int64, err error)
	WriteJSONTo(w io.Writer) (n int64, err error)
	WriteCLITo(w io.Writer) error
}

type JResponseReader interface {
	ReadXMLFrom(r io.Reader) (n int64, err error)
	ReadJSONFrom(r io.Reader) (n int64, err error)
}

type ResponseReaderWriter interface {
	JResponseReader
	JResponseWriter
}
