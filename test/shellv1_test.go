package test

import (
	"dadb"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func shellV1(t *testing.T, d dadb.Dadb) {
	withStream(t, d, "shell:echo hello", func(t *testing.T, d dadb.Dadb, stream dadb.Stream) {
		output, err := io.ReadAll(stream)
		require.NoError(t, err)
		require.Equal(t, string(output), "hello\n")
	})
}
