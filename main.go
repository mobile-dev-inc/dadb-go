package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func main() {
	host := "localhost"
	port := 5555

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	err = Connect(conn)
	if err != nil {
		panic(err)
	}

	conn.Close()
}

const AuthTypeToken = 1
const AuthTypeSignature = 2
const AuthTypeRsaPublic = 3

const CmdAuth = 0x48545541
const CmdCnxn = 0x4e584e43
const CmdOpen = 0x4e45504f
const CmdOkay = 0x59414b4f
const CmdClse = 0x45534c43
const CmdWrte = 0x45545257

const ConnectVersion = 0x01000000
const ConnectMaxData = 1024 * 1024

var ConnectPayload = []byte("host::\u0000")

type Packet struct {
	Command uint32
	Arg0    uint32
	Arg1    uint32
	Payload []byte
}

type PacketHeader struct {
	Command       uint32
	Arg0          uint32
	Arg1          uint32
	PayloadLength uint32
	Checksum      uint32
	Magic         uint32
}

func Connect(conn io.ReadWriter) error {
	err := WriteConnect(conn)
	if err != nil {
		return err
	}
	p, err := ReadPacket(conn)
	if err != nil {
		return err
	}
	fmt.Printf("==%v", p)
	return nil
}

func WriteConnect(w io.Writer) error {
	err := WritePacket(w, Packet{CmdCnxn, ConnectVersion, ConnectMaxData, ConnectPayload})
	if err != nil {
		return err
	}
	return nil
}

func ReadPacket(r io.Reader) (Packet, error) {
	var h PacketHeader
	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return Packet{}, err
	}

	payload := make([]byte, h.PayloadLength)
	_, err = r.Read(payload)
	if err != nil {
		return Packet{}, err
	}

	return Packet{
		Command: h.Command,
		Arg0:    h.Arg0,
		Arg1:    h.Arg1,
		Payload: payload,
	}, nil
}

func WritePacket(w io.Writer, p Packet) error {
	h := PacketHeader{
		Command:       p.Command,
		Arg0:          p.Arg0,
		Arg1:          p.Arg1,
		PayloadLength: uint32(len(p.Payload)),
		Checksum:      GetPayloadChecksum(p.Payload),
		Magic:         p.Command ^ 0xFFFFFFFF,
	}

	err := binary.Write(w, binary.LittleEndian, h)
	if err != nil {
		return err
	}

	if p.Payload != nil {
		_, err = w.Write(p.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetPayloadChecksum(payload []byte) uint32 {
	if payload == nil {
		return 0
	}
	var checksum uint32 = 0
	for i := 0; i < len(payload); i++ {
		checksum += uint32(payload[i])
	}
	return checksum
}
