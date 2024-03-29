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
		require.Equal(t, "hello\n", string(output))
	})
}
