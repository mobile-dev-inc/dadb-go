package test

import (
	"dadb"
	"dadb/adbd"
	"net"
	"testing"
)

func TestDadb(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:5555")
	if err != nil {
		t.Fatal(err)
	}
	dadb, err := adbd.Connect(conn)
	if err != nil {
		t.Fatal(err)
	}
	runDadbTest(t, dadb)
}

func runDadbTest(t *testing.T, d dadb.Dadb) {
	t.Run("Hello", func(t *testing.T) {

	})
}
