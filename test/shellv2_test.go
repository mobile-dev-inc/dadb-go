package test

import (
	"dadb"
	"testing"
)

func shellV2(t *testing.T, d dadb.Dadb) {
	requireShellOutput(t, d, "echo hello", "hello\n")
}
