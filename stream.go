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

func (s *Stream) Write(p []byte) (int, error) {
	// TODO what about when len(p) > s.connection.connectionResponse.maxPayloadSize?
	err := writePacket(s.connection.rw, packet{
		Command: cmdWrte,
		Arg0:    s.localId,
		Arg1:    s.remoteId,
		Payload: p,
	})

	if err != nil {
		return 0, err
	}

	pkt, err := s.readPacket()
	if err != nil {
		return 0, err
	}

	if pkt.Command != cmdOkay {
		return 0, fmt.Errorf("unexpected: command received 0x%x", pkt.Command)
	}

	return len(p), nil
}

func (s *Stream) getPayload() ([]byte, error) {
	if len(s.payload) > 0 {
		return s.payload, nil
	}

	pkt, err := s.readPacket()

	if err != nil {
		return nil, err
	}

	if pkt.Command == cmdClse {
		return nil, io.EOF
	}

	if pkt.Command != cmdWrte {
		return nil, fmt.Errorf("unexpected: command received 0x%x", pkt.Command)
	}

	err = writePacket(s.connection.rw, packet{
		Command: cmdOkay,
		Arg0:    s.localId,
		Arg1:    s.remoteId,
		Payload: nil,
	})

	if err != nil {
		return nil, err
	}

	return pkt.Payload, nil
}

func (s *Stream) readPacket() (packet, error) {
	ch := s.connection.getStreamChannel(s.localId)
	if ch == nil {
		return packet{}, fmt.Errorf("could not find channel for read: local id=0x%x", s.localId)
	}

	return <-ch, nil
}
