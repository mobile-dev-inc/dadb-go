package adbd

import (
	"dadb"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type connection struct {
	sync.RWMutex
	io.Reader
	io.Writer
	io.Closer
	connectionResponse connectionResponse

	nextLocalId   uint32
	channels      map[uint32]map[uint32]chan packet
	closedStreams map[uint32]struct{}
}

func Connect(conn net.Conn) (dadb.Dadb, error) {
	response, err := connect(conn)
	if err != nil {
		return nil, err
	}

	connection := connection{
		Reader:             conn,
		Writer:             conn,
		Closer:             conn,
		connectionResponse: response,
		nextLocalId:        0,
		channels:           make(map[uint32]map[uint32]chan packet),
		closedStreams:      make(map[uint32]struct{}),
	}

	go func() {
		for {
			p, err := readPacket(conn)
			if err != nil {
				log.Printf("TODO: Error in AdbdConnection goroutine: %v\n", err)
				return
			}

			localId := p.Arg1
			cmd := p.Command
			if cmd == cmdClse {
				connection.closeStream(localId)
			} else {
				ch := connection.getChannel(localId, cmd)
				// No need to lock since we only close channels from this goroutine. Also, the channel
				// shouldn't be closed at this point based on the adb protocol, but we check just in case
				// to avoid a panic.
				_, closed := connection.closedStreams[localId]
				if !closed {
					ch <- p
				}
			}
		}
	}()

	return &connection, nil
}

func (c *connection) Open(destination string) (dadb.Stream, error) {
	localId := atomic.AddUint32(&c.nextLocalId, 1)

	err := writeOpen(c, localId, destination)
	if err != nil {
		return nil, err
	}

	p := <-c.getChannel(localId, cmdOkay)
	remoteId := p.Arg0

	return &stream{
		connection: c,
		localId:    localId,
		remoteId:   remoteId,
	}, nil
}

func (c *connection) SupportsFeature(feature string) bool {
	_, ok := c.connectionResponse.features[feature]
	return ok
}

func (c *connection) closeStream(localId uint32) {
	c.Lock()
	defer c.Unlock()

	_, alreadyClosed := c.closedStreams[localId]
	if alreadyClosed {
		return
	}

	c.closedStreams[localId] = struct{}{}
	for _, ch := range c.channels[localId] {
		close(ch)
	}
}

func (c *connection) getChannel(localId uint32, cmd uint32) chan packet {
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
		_, closed := c.closedStreams[localId]
		if closed {
			close(ch)
		}
	}
	return ch
}
