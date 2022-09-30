package test

import (
	"dadb"
	"dadb/adbd"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func shellV1(t *testing.T, d dadb.Dadb) {
	stream, err := d.Open("shell:echo hello")
	require.Nil(t, err)
	output, err := io.ReadAll(stream)
	require.Nil(t, err)
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
	runDadbTest(t, adbdDadb, "adbd")
}

func createAdbdDadb(t *testing.T) dadb.Dadb {
	conn, err := net.Dial("tcp", "localhost:5555")
	require.Nil(t, err)
	dadb, err := adbd.Connect(conn)
	require.Nil(t, err)
	return dadb
}
