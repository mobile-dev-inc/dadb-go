package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	host := "localhost"
	port := 5555

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	err, connectionResponse := Connect(conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", connectionResponse)

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

type ConnectionResponse struct {
	version        uint32
	maxPayloadSize uint32
	props          map[string]string
	features       map[string]struct{}
}

func Connect(conn io.ReadWriter) (error, ConnectionResponse) {
	err := WriteConnect(conn)
	if err != nil {
		return err, ConnectionResponse{}
	}

	p, err := ReadPacket(conn)
	if err != nil {
		return err, ConnectionResponse{}
	}

	if p.Command == CmdAuth {
		panic("Not Implemented")
	}

	if p.Command != CmdCnxn {
		return fmt.Errorf("connection failed: unexpected command 0x%x", p.Command), ConnectionResponse{}
	}

	err, connectionResponse := ParseConnectionResponse(p)
	if err != nil {
		return err, ConnectionResponse{}
	}

	return nil, connectionResponse
}

// ParseConnectionResponse
// eg. device::ro.product.name=sdk_phone_arm64;ro.product.model=Android SDK built for arm64;ro.product.device=emulator_arm64;features=sendrecv_v2_brotli,remount_shell,sendrecv_v2,abb_exec,fixed_push_mkdir,fixed_push_symlink_timestamp,abb,shell_v2,cmd,ls_v2,apex,stat_v2
func ParseConnectionResponse(p Packet) (error, ConnectionResponse) {
	connectionStr := string(p.Payload)
	propsString := strings.SplitN(connectionStr, "device::", 2)[1]

	props := make(map[string]string)
	for _, prop := range strings.Split(propsString, ";") {
		propParts := strings.SplitN(prop, "=", 2)
		props[propParts[0]] = propParts[1]
	}

	featuresString, exists := props["features"]
	if !exists {
		return fmt.Errorf("failed to parse connection string: features not found (%s)", connectionStr), ConnectionResponse{}
	}

	features := make(map[string]struct{})
	for _, feature := range strings.Split(featuresString, ",") {
		features[feature] = struct{}{}
	}

	return nil, ConnectionResponse{
		version:        p.Arg0,
		maxPayloadSize: p.Arg1,
		props:          props,
		features:       features,
	}
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
