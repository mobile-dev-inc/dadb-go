package dadb

import (
	"dadb/adbd"
	"dadb/adbserver"
	"io"
	"net"
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

func CreateAdbd(conn net.Conn) (Dadb, error) {
	connection, err := adbd.Connect(conn)
	if err != nil {
		return Dadb{}, err
	}
	return Dadb{Connection: connection}, nil
}

func CreateAdbServer(address string, deviceQuery string) (Dadb, error) {
	connection, err := adbserver.Connect(address, deviceQuery)
	if err != nil {
		return Dadb{}, err
	}
	return Dadb{Connection: connection}, nil
}
