package main

import (
	"fmt"
	"net"
)

func main() {
	host := "localhost"
	port := 5555

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	err, connectionResponse := Connect(conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", connectionResponse)

	conn.Close()
}
