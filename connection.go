package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Connection struct {
	sync.RWMutex

	rw                 io.ReadWriter
	closer             io.Closer
	connectionResponse connectionResponse

	nextLocalId    uint32
	streamChannels map[uint32]chan packet
}

func Connect(conn net.Conn) (*Connection, error) {
	response, err := connect(conn)
	if err != nil {
		return nil, err
	}

	connection := Connection{
		rw:                 conn,
		closer:             conn,
		connectionResponse: response,
		nextLocalId:        0,
		streamChannels:     make(map[uint32]chan packet),
	}

	go func() {
		for {
			p, err := readPacket(conn)
			if err != nil {
				log.Printf("TODO: Error in Connection goroutine: %v\n", err)
				return
			}

			localId := p.Arg1
			ch := connection.getStreamChannel(localId)
			if ch == nil {
				log.Printf("TODO: Error in Connection goroutine: no channel for localId 0x%x\n", localId)
				return
			}
			ch <- p
		}
	}()

	return &connection, nil
}

func (c *Connection) Open(destination string) (*Stream, error) {
	localId := atomic.AddUint32(&c.nextLocalId, 1)

	ch := make(chan packet, 100)

	c.Lock()
	c.streamChannels[localId] = ch
	err := writeOpen(c.rw, localId, destination)
	if err != nil {
		return nil, err
	}
	c.Unlock()

	p := <-ch
	if p.Command != cmdOkay {
		return nil, fmt.Errorf("unexpected command: 0x%x", p.Arg0)
	}

	remoteId := p.Arg0

	return &Stream{
		connection: c,
		localId:    localId,
		remoteId:   remoteId,
		ch:         ch,
	}, nil
}

func (c *Connection) getStreamChannel(localId uint32) chan packet {
	c.RLock()
	defer c.RUnlock()
	return c.streamChannels[localId]
}
