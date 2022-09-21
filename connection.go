package main

import (
	"io"
	"net"
	"sync/atomic"
)

type Connection struct {
	closer             io.Closer
	rw                 io.ReadWriter
	connectionResponse connectionResponse
	nextLocalId        uint32
}

type Stream struct {
	connection *Connection
	localId    uint32
	remoteId   uint32
}

func Connect(conn net.Conn) (error, *Connection) {
	err, response := connect(conn)
	if err != nil {
		return err, nil
	}
	return nil, &Connection{
		closer:             conn,
		rw:                 conn,
		connectionResponse: response,
		nextLocalId:        0,
	}
}

func (c *Connection) Open(destination string) (error, *Stream) {
	localId := atomic.AddUint32(&c.nextLocalId, 1)

	err := writeOpen(c.rw, localId, destination)
	if err != nil {
		return err, nil
	}

	panic(0)
}
