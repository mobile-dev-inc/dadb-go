package adbd

import (
	"dadb"
	"net"
)

func CreateDadb(conn net.Conn) (dadb.Dadb, error) {
	connection, err := connect(conn)
	if err != nil {
		return dadb.Dadb{}, err
	}
	return dadb.Dadb{Connection: &connection}, nil
}
