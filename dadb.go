package dadb

import (
	"fmt"
	"io"
)

type Connection interface {
	Open(destination string) (Stream, error)
	SupportsFeature(feature string) bool
}

type Stream interface {
	io.Reader
	io.Writer
	io.Closer
}

type Dadb struct {
	Connection
}

func (d Dadb) Shell(command string) (ShellResponse, error) {
	stream, err := d.OpenShell(command)
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

func (d Dadb) OpenShell(command string) (ShellStream, error) {
	stream, err := d.Open(fmt.Sprintf("shell,v2,raw:%s", command))
	if err != nil {
		return ShellStream{}, err
	}
	return ShellStream{s: stream}, nil
}
