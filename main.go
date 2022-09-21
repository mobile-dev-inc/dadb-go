package main

import (
	"io"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5555")
	if err != nil {
		panic(err)
	}

	connection, err := Connect(conn)
	if err != nil {
		panic(err)
	}

	stream, err := connection.Open("shell:echo hello")
	if err != nil {
		panic(err)
	}

	all, err := io.ReadAll(stream)
	if err != nil {
		panic(err)
	}

	log.Println(string(all))
}
