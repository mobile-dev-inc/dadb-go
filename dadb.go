package dadb

import (
	"fmt"
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

func Shell(d Dadb, command string) (ShellResponse, error) {
	stream, err := OpenShell(d, command)
	if err != nil {
		return ShellResponse{}, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer stream.Close()

	if err != nil {
		return ShellResponse{}, err
	}

	return stream.ReadAll()
}

func OpenShell(d Dadb, command string) (ShellStream, error) {
	stream, err := d.Open(fmt.Sprintf("shell,v2,raw:%s", command))
	if err != nil {
		return ShellStream{}, err
	}
	return ShellStream{s: stream}, nil
}
