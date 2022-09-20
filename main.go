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

	Connect(conn)

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
	command uint32
	arg0    uint32
	arg1    uint32
	payload []byte
}

func Connect(conn io.ReadWriter) error {
	err := WriteConnect(conn)
	if err != nil {
		return err
	}
	return nil
}

func WriteConnect(w io.Writer) error {
	err := WritePacket(w, &Packet{CmdCnxn, ConnectVersion, ConnectMaxData, ConnectPayload})
	if err != nil {
		return err
	}
	return nil
}

func WritePacket(
	w io.Writer,
	p *Packet,
) error {
	err := WriteLe(w, p.command)
	if err != nil {
		return err
	}

	err = WriteLe(w, p.arg0)
	if err != nil {
		return err
	}

	err = WriteLe(w, p.arg1)
	if err != nil {
		return err
	}

	if p.payload == nil {
		err = WriteLe(w, 0)
		if err != nil {
			return err
		}
		err = WriteLe(w, 0)
		if err != nil {
			return err
		}
	} else {
		err = WriteLe(w, uint32(len(p.payload)))
		if err != nil {
			return err
		}
		checksum := GetPayloadChecksum(p.payload)
		err = WriteLe(w, checksum)
		if err != nil {
			return err
		}
	}

	err = WriteLe(w, p.command^0xFFFFFFFF)
	if err != nil {
		return err
	}

	if p.payload != nil {
		_, err = w.Write(p.payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetPayloadChecksum(payload []byte) uint32 {
	var checksum uint32 = 0
	for i := 0; i < len(payload); i++ {
		checksum += uint32(payload[i])
	}
	return checksum
}

func WriteLe(w io.Writer, i uint32) error {
	return binary.Write(w, binary.LittleEndian, i)
}
