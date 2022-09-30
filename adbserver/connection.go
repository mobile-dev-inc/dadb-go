package adbserver

import (
	"dadb"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type connection struct {
	address     string
	deviceQuery string
	features    map[string]struct{}
}

type stream struct {
	io.Reader
	io.Writer
	io.Closer
}

func Connect(address string, deviceQuery string) (dadb.Dadb, error) {
	features, err := readFeatures(address, deviceQuery)
	if err != nil {
		return nil, err
	}
	return connection{
		address:     address,
		deviceQuery: deviceQuery,
		features:    features,
	}, nil
}

func (c connection) Open(destination string) (dadb.Stream, error) {
	conn, err := open(c.address, c.deviceQuery, destination)
	if err != nil {
		return nil, err
	}
	return stream{
		Reader: conn,
		Writer: conn,
		Closer: conn,
	}, nil
}

func (c connection) SupportsFeature(feature string) bool {
	_, ok := c.features[feature]
	return ok
}

func readFeatures(address string, deviceQuery string) (map[string]struct{}, error) {
	rw, err := open(address, deviceQuery, "host:features")
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(rw)
	if err != nil {
		return nil, err
	}
	features := make(map[string]struct{})
	for _, feature := range strings.Split(string(bytes), ",") {
		features[feature] = struct{}{}
	}
	return features, nil
}

func open(address string, deviceQuery string, destination string) (net.Conn, error) {
	// TODO: Ensure server is running
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	err = send(conn, deviceQuery)
	if err != nil {
		return nil, err
	}
	err = send(conn, destination)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func send(rw io.ReadWriter, command string) error {
	err := writeString(rw, command)
	if err != nil {
		return err
	}
	response, err := readStringLen(rw, 4)
	if err != nil {
		return err
	}
	if response != "OKAY" {
		errMsg, err := readString(rw)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to send command (%s): %s", command, errMsg)
	}
	return nil
}

func writeString(w io.Writer, string string) error {
	_, err := fmt.Fprintf(w, "%04x", len(string))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(string))
	if err != nil {
		return err
	}
	return nil
}

func readString(r io.Reader) (string, error) {
	encodedLength, err := readStringLen(r, 4)
	if err != nil {
		return "", err
	}
	length, err := strconv.ParseInt(encodedLength, 16, 32)
	if err != nil {
		return "", err
	}
	return readStringLen(r, int(length))
}

func readStringLen(r io.Reader, len int) (string, error) {
	bytes := make([]byte, len)
	_, err := io.ReadFull(r, bytes)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
