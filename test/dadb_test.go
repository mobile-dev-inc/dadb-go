package test

import (
	"dadb"
	"dadb/adbd"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func TestDadb(t *testing.T) {
	adbdDadb := createAdbdDadb(t)
	runDadbTest(t, adbdDadb)
}

func createAdbdDadb(t *testing.T) dadb.Dadb {
	conn, err := net.Dial("tcp", "localhost:5555")
	require.Nil(t, err)
	dadb, err := adbd.Connect(conn)
	require.Nil(t, err)
	return dadb
}

func runDadbTest(t *testing.T, d dadb.Dadb) {
	t.Run("shellV1", func(t *testing.T) {
		stream, err := d.Open("shell:echo hello")
		require.Nil(t, err)
		output, err := io.ReadAll(stream)
		require.Nil(t, err)
		require.Equal(t, string(output), "hello\n")
	})
}
