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
	Id      byte
	Payload []byte
}

type ShellResponse struct {
	Output      string
	ErrorOutput string
	ExitCode    int
}

type packetHeader struct {
	Id  byte
	Len uint32
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
		switch packet.Id {
		case IdExit:
			return ShellResponse{
				Output:      string(output),
				ErrorOutput: string(errorOutput),
				ExitCode:    int(packet.Payload[0]),
			}, nil
		case IdStdout:
			output = append(output, packet.Payload...)
		case IdStderr:
			errorOutput = append(errorOutput, packet.Payload...)
		}
	}
}

func (s ShellStream) Read() (ShellPacket, error) {
	header := packetHeader{}
	err := binary.Read(s.s, binary.LittleEndian, &header)
	if err != nil {
		return ShellPacket{}, err
	}
	payload := make([]byte, header.Len)
	_, err = s.s.Read(payload)
	if err != nil {
		return ShellPacket{}, err
	}
	return ShellPacket{
		Id:      header.Id,
		Payload: payload,
	}, nil
}

func (s ShellStream) WriteString(string string) error {
	return s.Write(IdStdin, []byte(string))
}

func (s ShellStream) Write(id byte, payload []byte) error {
	err := binary.Write(s.s, binary.LittleEndian, packetHeader{
		Id:  id,
		Len: uint32(len(payload)),
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

func (s ShellStream) Close() error {
	err := s.s.Close()
	if err != nil {
		return err
	}
	return nil
}
