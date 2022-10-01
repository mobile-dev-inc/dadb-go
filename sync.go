package dadb

import (
	"encoding/binary"
	"io"
)

const list = "LIST"
const recv = "RECV"
const send = "SEND"
const stat = "STAT"
const data = "DATA"
const done = "DONE"
const okay = "OKAY"
const quit = "QUIT"
const fail = "FAIL"

type SyncStream struct {
	s Stream
}

type syncPacket struct {
	id  string
	arg uint32
}

func (s *SyncStream) read() (syncPacket, error) {
	idBytes := make([]byte, 4)
	_, err := io.ReadFull(s.s, idBytes)
	if err != nil {
		return syncPacket{}, err
	}

	argBytes := make([]byte, 4)
	arg := binary.LittleEndian.Uint32(argBytes)

	return syncPacket{
		id:  string(idBytes),
		arg: arg,
	}, nil
}
