package test

import (
	"dadb"
	"github.com/stretchr/testify/require"
	"testing"
)

func shellV2(t *testing.T, d dadb.Dadb) {
	response, err := d.Shell("echo hello")
	require.NoError(t, err)
	require.Equal(t, "hello\n", response.Output)
}
