package test

import (
	"bytes"
	"dadb"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func push(t *testing.T, d dadb.Dadb) {
	r := bytes.NewReader([]byte("Hello World!"))
	err := dadb.Push(d, r, remoteFilePath, 0o666, 0)
	require.NoError(t, err)
	requireShellOutput(t, d, fmt.Sprintf("cat %s", remoteFilePath), "Hello World!")
}

func pull(t *testing.T, d dadb.Dadb) {
	requireShell(t, d, fmt.Sprintf("echo Hello World! > %s", remoteFilePath))

	var b bytes.Buffer
	err := dadb.Pull(d, &b, remoteFilePath)
	require.NoError(t, err)

	assert.Equal(t, "Hello World!\n", b.String())
}
