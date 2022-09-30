package main

import "io"

type Dadb interface {
	Open(destination string) (Stream, error)
}

type Stream interface {
	io.Reader
	io.Writer

	SupportsFeature(feature string) bool
}
