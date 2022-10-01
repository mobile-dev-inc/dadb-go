package dadb

import (
	"encoding/binary"
)

const IdStdin byte = 0
const IdStdout byte = 1
const IdStderr byte = 2
const IdExit byte = 3
const IdCloseStdin byte = 3

type ShellPacket struct {
	id      byte
	payload []byte
}

type ShellResponse struct {
	output      string
	errorOutput string
	exitCode    int
}

type packetHeader struct {
	id  byte
	len uint32
}

type ShellStream struct {
	s Stream
}

func (s ShellStream) ReadAll() (ShellResponse, error) {
	output := make([]byte, 0)
	errorOutput := make([]byte, 0)

	for {
		packet, err := s.Read()
		if err != nil {
			return ShellResponse{}, err
		}
		switch packet.id {
		case IdExit:
			return ShellResponse{
				output:      string(output),
				errorOutput: string(errorOutput),
				exitCode:    int(packet.payload[0]),
			}, nil
		case IdStdin:
			output = append(output, packet.payload...)
			break
		case IdStderr:
			errorOutput = append(errorOutput, packet.payload...)
			break
		}
	}
}

func (s ShellStream) Read() (ShellPacket, error) {
	header := packetHeader{}
	err := binary.Read(s.s, binary.LittleEndian, &header)
	if err != nil {
		return ShellPacket{}, err
	}
	payload := make([]byte, header.len)
	_, err = s.s.Read(payload)
	if err != nil {
		return ShellPacket{}, err
	}
	return ShellPacket{
		id:      header.id,
		payload: payload,
	}, nil
}

func (s ShellStream) WriteString(string string) error {
	return s.Write(IdStdin, []byte(string))
}

func (s ShellStream) Write(id byte, payload []byte) error {
	err := binary.Write(s.s, binary.LittleEndian, packetHeader{
		id:  id,
		len: uint32(len(payload)),
	})
	if err != nil {
		return err
	}
	_, err = s.s.Write(payload)
	if err != nil {
		return err
	}
	return nil
}
