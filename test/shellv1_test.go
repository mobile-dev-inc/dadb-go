package test

import (
	"dadb"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func shellV1(t *testing.T, d dadb.Dadb) {
	withStream(t, d, "shell:echo hello", func(stream dadb.Stream) {
		output, err := io.ReadAll(stream)
		require.NoError(t, err)
		require.Equal(t, string(output), "hello\n")
	})
}

func shellV1Close(t *testing.T, d dadb.Dadb) {
	stream, err := d.Open("shell:")
	require.NoError(t, err)
	err = stream.Close()
	require.NoError(t, err)
	read, err := stream.Read(make([]byte, 10))
	require.Equal(t, read, 0)
	require.Error(t, err)
}
