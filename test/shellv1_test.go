package test

import (
	"dadb"
	"github.com/stretchr/testify/require"
	"io"
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
