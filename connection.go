package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Connection struct {
	close              chan struct{}
	conn               *Connection
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

	connection := Connection{
		close:              make(chan struct{}),
		rw:                 conn,
		connectionResponse: response,
		nextLocalId:        0,
	}

	go connection.loop()

	return nil, &connection
}

func (c *Connection) loop() {
	for {
		select {
		case <-c.close:
			err := c.conn.Close()
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
}

func (c *Connection) Close() error {
	c.close <- struct{}{}
	return nil
}

func (c *Connection) Open(destination string) (error, *Stream) {
	localId := atomic.AddUint32(&c.nextLocalId, 1)

	err := writeOpen(c.rw, localId, destination)
	if err != nil {
		return err, nil
	}

	panic(0)
}

func tmp() {
	host := "localhost"
	port := 5555

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	err, connectionResponse := connect(conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", connectionResponse)

	err = writeOpen(conn, 1, "shell:echo hello")
	if err != nil {
		panic(err)
	}

	packet, err := readPacket(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println(packet)

	packet, _ = readPacket(conn)
	fmt.Println(string(packet.Payload))

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
