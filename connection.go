package main

import (
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

	nextLocalId uint32
	channels    map[uint32]map[uint32]chan packet
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
		channels:           make(map[uint32]map[uint32]chan packet),
	}

	go func() {
		for {
			p, err := readPacket(conn)
			if err != nil {
				log.Printf("TODO: Error in Connection goroutine: %v\n", err)
				return
			}

			localId := p.Arg1
			cmd := p.Command
			ch := connection.getChannel(localId, cmd)
			ch <- p
		}
	}()

	return &connection, nil
}

func (c *Connection) Open(destination string) (*Stream, error) {
	localId := atomic.AddUint32(&c.nextLocalId, 1)

	err := writeOpen(c.rw, localId, destination)
	if err != nil {
		return nil, err
	}

	p := <-c.getChannel(localId, cmdOkay)
	remoteId := p.Arg0

	return &Stream{
		connection: c,
		localId:    localId,
		remoteId:   remoteId,
	}, nil
}

func (c *Connection) getChannel(localId uint32, cmd uint32) chan packet {
	// Fast path: Channel already exists - Only acquire read lock
	c.RLock()
	m := c.channels[localId]
	if m != nil {
		ch := m[cmd]
		if ch != nil {
			c.RUnlock()
			return ch
		}
	}
	c.RUnlock()

	// Slow path: Channel does not exist - Acquire write lock
	c.Lock()
	defer c.Unlock()

	m = c.channels[localId]
	if m == nil {
		m = make(map[uint32]chan packet)
		c.channels[localId] = m
	}
	ch := m[cmd]
	if ch == nil {
		ch = make(chan packet, 100)
		m[cmd] = ch
	}
	return ch
}
