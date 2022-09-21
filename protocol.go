package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const authTypeToken = 1
const authTypeSignature = 2
const authTypeRsaPublic = 3

const cmdAuth = 0x48545541
const cmdCnxn = 0x4e584e43
const cmdOpen = 0x4e45504f
const cmdOkay = 0x59414b4f
const cmdClse = 0x45534c43
const cmdWrte = 0x45545257

const connectVersion = 0x01000000
const connectMaxData = 1024 * 1024

var connectPayload = []byte("host::\u0000")

type packet struct {
	Command uint32
	Arg0    uint32
	Arg1    uint32
	Payload []byte
}

type packetHeader struct {
	Command       uint32
	Arg0          uint32
	Arg1          uint32
	PayloadLength uint32
	Checksum      uint32
	Magic         uint32
}

type connectionResponse struct {
	version        uint32
	maxPayloadSize uint32
	props          map[string]string
	features       map[string]struct{}
}

func writeOpen(w io.Writer, localId uint32, destination string) error {
	destinationBytes := append([]byte(destination), 0)
	return writePacket(w, packet{
		Command: cmdOpen,
		Arg0:    localId,
		Arg1:    0,
		Payload: destinationBytes,
	})
}

func connect(conn io.ReadWriter) (error, connectionResponse) {
	err := writeConnect(conn)
	if err != nil {
		return err, connectionResponse{}
	}

	p, err := readPacket(conn)
	if err != nil {
		return err, connectionResponse{}
	}

	if p.Command == cmdAuth {
		panic("Not Implemented")
	}

	if p.Command != cmdCnxn {
		return fmt.Errorf("connection failed: unexpected command 0x%x", p.Command), connectionResponse{}
	}

	err, response := parseConnectionResponse(p)
	if err != nil {
		return err, connectionResponse{}
	}

	return nil, response
}

// parseConnectionResponse
// eg. device::ro.product.name=sdk_phone_arm64;ro.product.model=Android SDK built for arm64;ro.product.device=emulator_arm64;features=sendrecv_v2_brotli,remount_shell,sendrecv_v2,abb_exec,fixed_push_mkdir,fixed_push_symlink_timestamp,abb,shell_v2,cmd,ls_v2,apex,stat_v2
func parseConnectionResponse(p packet) (error, connectionResponse) {
	connectionStr := string(p.Payload)
	propsString := strings.SplitN(connectionStr, "device::", 2)[1]

	props := make(map[string]string)
	for _, prop := range strings.Split(propsString, ";") {
		propParts := strings.SplitN(prop, "=", 2)
		props[propParts[0]] = propParts[1]
	}

	featuresString, exists := props["features"]
	if !exists {
		return fmt.Errorf("failed to parse connection string: features not found (%s)", connectionStr), connectionResponse{}
	}

	features := make(map[string]struct{})
	for _, feature := range strings.Split(featuresString, ",") {
		features[feature] = struct{}{}
	}

	return nil, connectionResponse{
		version:        p.Arg0,
		maxPayloadSize: p.Arg1,
		props:          props,
		features:       features,
	}
}

func writeConnect(w io.Writer) error {
	err := writePacket(w, packet{cmdCnxn, connectVersion, connectMaxData, connectPayload})
	if err != nil {
		return err
	}
	return nil
}

func readPacket(r io.Reader) (packet, error) {
	var h packetHeader
	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return packet{}, err
	}

	payload := make([]byte, h.PayloadLength)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return packet{}, err
	}

	return packet{
		Command: h.Command,
		Arg0:    h.Arg0,
		Arg1:    h.Arg1,
		Payload: payload,
	}, nil
}

func writePacket(w io.Writer, p packet) error {
	h := packetHeader{
		Command:       p.Command,
		Arg0:          p.Arg0,
		Arg1:          p.Arg1,
		PayloadLength: uint32(len(p.Payload)),
		Checksum:      getPayloadChecksum(p.Payload),
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

func getPayloadChecksum(payload []byte) uint32 {
	if payload == nil {
		return 0
	}
	var checksum uint32 = 0
	for i := 0; i < len(payload); i++ {
		checksum += uint32(payload[i])
	}
	return checksum
}
