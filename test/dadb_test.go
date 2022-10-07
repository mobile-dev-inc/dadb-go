package test

import (
	"dadb"
	"dadb/adbd"
	"dadb/adbserver"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

const remoteFilePath = "/data/local/tmp/testfile"

func withStream(
	t *testing.T,
	c dadb.Dadb,
	destination string,
	f func(stream dadb.Stream),
) {
	stream, err := c.Open(destination)
	require.NoError(t, err)
	defer close(t, stream)
	f(stream)
}

func runDadbTest(t *testing.T, d dadb.Dadb, prefix string) {
	run := func(name string, f func(t *testing.T, d dadb.Dadb)) {
		_, err := dadb.Shell(d, fmt.Sprintf("rm -rf %s", remoteFilePath))
		if err != nil {
			require.NoError(t, err)
		}
		testName := fmt.Sprintf("%s/%s", prefix, name)
		t.Run(testName, func(t *testing.T) {
			f(t, d)
		})
	}
	run("shellV1", shellV1)
	run("shellV2", shellV2)
	run("push", push)
}

func TestDadb(t *testing.T) {
	adbdDadb := connectAdbd(t)
	adbServerDadb := connectAdbServer(t)
	runDadbTest(t, adbdDadb, "adbd")
	runDadbTest(t, adbServerDadb, "adbserver")
}

func requireShell(t *testing.T, d dadb.Dadb, command string) string {
	response, err := dadb.Shell(d, command)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, response.ExitCode, 0)
	return response.Output
}

func requireShellOutput(t *testing.T, d dadb.Dadb, command string, expected string) {
	output := requireShell(t, d, command)
	assert.Equal(t, expected, output)
}

func connectAdbd(t *testing.T) dadb.Dadb {
	c, err := adbd.Connect("tcp", "localhost:5555")
	require.NoError(t, err)
	return c
}

func connectAdbServer(t *testing.T) dadb.Dadb {
	c, err := adbserver.Connect("localhost:5037", "host:transport-any")
	require.NoError(t, err)
	return c
}

func close(t *testing.T, c io.Closer) {
	err := c.Close()
	require.NoError(t, err)
}
