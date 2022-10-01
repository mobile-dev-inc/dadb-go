package test

import (
	"dadb"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func withStream(
	t *testing.T,
	c dadb.Connection,
	destination string,
	f func(stream dadb.Stream),
) {
	stream, err := c.Open(destination)
	require.NoError(t, err)
	defer close(t, stream)
	f(stream)
}

func runDadbTest(t *testing.T, d dadb.Dadb, prefix string) {
	run := func(name string, f func(t *testing.T, d dadb.Connection)) {
		testName := fmt.Sprintf("%s/%s", prefix, name)
		t.Run(testName, func(t *testing.T) {
			f(t, d)
		})
	}
	run("shellV1", shellV1)
}

func TestDadb(t *testing.T) {
	adbdDadb := connectAdbd(t)
	adbServerDadb := connectAdbServer(t)
	runDadbTest(t, adbdDadb, "adbd")
	runDadbTest(t, adbServerDadb, "adbserver")
}

func connectAdbd(t *testing.T) dadb.Dadb {
	conn, err := net.Dial("tcp", "localhost:5555")
	require.NoError(t, err)
	c, err := dadb.CreateAdbd(conn)
	require.NoError(t, err)
	return c
}

func connectAdbServer(t *testing.T) dadb.Dadb {
	c, err := dadb.CreateAdbServer("localhost:5037", "host:transport-any")
	require.NoError(t, err)
	return c
}

func close(t *testing.T, c io.Closer) {
	err := c.Close()
	require.NoError(t, err)
}
