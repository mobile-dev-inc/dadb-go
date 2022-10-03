package dadb

import (
	"io"
)

type Dadb interface {
	Open(destination string) (Stream, error)
	SupportsFeature(feature string) bool
}

type Stream interface {
	io.Reader
	io.Writer
	io.Closer
}
