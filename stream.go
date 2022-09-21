package main

import (
	"fmt"
	"io"
)

type Stream struct {
	connection *Connection
	localId    uint32
	remoteId   uint32

	ch      chan packet
	payload []byte
}

func (s *Stream) Read(p []byte) (int, error) {
	payload, err := s.getPayload()
	if err != nil {
		return 0, err
	}

	n := copy(p, payload)

	s.payload = payload[n:]

	return n, nil
}

func (s *Stream) getPayload() ([]byte, error) {
	if len(s.payload) > 0 {
		return s.payload, nil
	}

	ch := s.connection.getStreamChannel(s.localId)
	if ch == nil {
		return nil, fmt.Errorf("could not find channel for read: local id=0x%x", s.localId)
	}

	p := <-ch

	if p.Command == cmdClse {
		return nil, io.EOF
	}

	if p.Command != cmdWrte {
		return nil, fmt.Errorf("unexpected: command received 0x%x", p.Command)
	}

	return p.Payload, nil
}
