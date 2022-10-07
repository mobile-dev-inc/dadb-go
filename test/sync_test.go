package test

import (
	"bytes"
	"dadb"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func push(t *testing.T, d dadb.Dadb) {
	r := bytes.NewReader([]byte("Hello World!"))
	err := dadb.Push(d, r, remoteFilePath, 0o666, 0)
	if err != nil {
		require.NoError(t, err)
	}
	requireShellOutput(t, d, fmt.Sprintf("cat %s", remoteFilePath), "Hello World!")
}
