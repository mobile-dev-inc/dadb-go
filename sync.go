package dadb

import (
	"encoding/binary"
	"errors"
	"fmt"
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

type syncPacket struct {
	id  string
	arg uint32
}

type syncStream struct {
	s       Stream
	payload []byte
}

func Push(dadb Dadb, r io.Reader, remotePath string, mode uint32, lastModifiedSec uint32) error {
	stream, err := dadb.Open("sync:")
	if err != nil {
		return err
	}

	remote := fmt.Sprintf("%s,%d", remotePath, mode)
	err = writeSyncPacketWithPayload(stream, send, []byte(remote))
	if err != nil {
		return err
	}

	syncStream := syncStream{s: stream}

	// If needed, we can try increasing the buffer size here to improve performance
	_, err = io.Copy(syncStream, r)
	if err != nil {
		return err
	}

	return writeSyncPacket(stream, done, lastModifiedSec)
}

func (s syncStream) Read(p []byte) (int, error) {
	if len(s.payload) == 0 {
		packet, err := readSyncPacket(s.s)
		if err != nil {
			return 0, err
		}
		switch packet.id {
		case done:
			return 0, io.EOF
		case fail:
			bytes, err := readNBytes(s.s, int(packet.arg))
			if err != nil {
				return 0, err
			}
			return 0, errors.New(string(bytes))
		case data:
			bytes, err := readNBytes(s.s, int(packet.arg))
			if err != nil {
				return 0, err
			}
			s.payload = bytes
		}
	}

	n := copy(p, s.payload)

	s.payload = s.payload[n:]

	return n, nil
}

func (s syncStream) Write(p []byte) (int, error) {
	err := writeSyncPacket(s.s, data, uint32(len(p)))
	if err != nil {
		return 0, err
	}
	return s.s.Write(p)
}

func readSyncPacket(r io.Reader) (syncPacket, error) {
	idBytes := make([]byte, 4)
	_, err := io.ReadFull(r, idBytes)
	if err != nil {
		return syncPacket{}, err
	}

	argBytes, err := readNBytes(r, 4)
	if err != nil {
		return syncPacket{}, err
	}
	arg := binary.LittleEndian.Uint32(argBytes)

	return syncPacket{
		id:  string(idBytes),
		arg: arg,
	}, nil
}

func readNBytes(r io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := r.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func writeSyncPacketWithPayload(w io.Writer, id string, payload []byte) error {
	err := writeSyncPacket(w, id, uint32(len(payload)))
	if err != nil {
		return err
	}
	_, err = w.Write(payload)
	if err != nil {
		return err
	}
	return nil
}

func writeSyncPacket(w io.Writer, id string, arg uint32) error {
	_, err := io.WriteString(w, id)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, arg)
	if err != nil {
		return err
	}

	return nil
}
