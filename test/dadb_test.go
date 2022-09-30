package test

import (
	"dadb"
	"dadb/adbd"
	"dadb/adbserver"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func shellV1(t *testing.T, d dadb.Dadb) {
	stream, err := d.Open("shell:echo hello")
	require.NoError(t, err)
	defer close(t, stream)
	output, err := io.ReadAll(stream)
	require.NoError(t, err)
	require.Equal(t, string(output), "hello\n")
}

func runDadbTest(t *testing.T, d dadb.Dadb, prefix string) {
	run := func(name string, f func(t *testing.T, d dadb.Dadb)) {
		testName := fmt.Sprintf("%s/%s", prefix, name)
		t.Run(testName, func(t *testing.T) {
			f(t, d)
		})
	}
	run("shellV1", shellV1)
}

func TestDadb(t *testing.T) {
	adbdDadb := createAdbdDadb(t)
	adbServerDadb := createAdbServerDadb(t)
	runDadbTest(t, adbdDadb, "adbd")
	runDadbTest(t, adbServerDadb, "adbserver")
}

func createAdbdDadb(t *testing.T) dadb.Dadb {
	conn, err := net.Dial("tcp", "localhost:5555")
	require.NoError(t, err)
	dadb, err := adbd.Connect(conn)
	require.NoError(t, err)
	return dadb
}

func createAdbServerDadb(t *testing.T) dadb.Dadb {
	dadb, err := adbserver.Connect("localhost:5037", "host:transport-any")
	require.NoError(t, err)
	return dadb
}

func close(t *testing.T, c io.Closer) {
	err := c.Close()
	require.NoError(t, err)
}
