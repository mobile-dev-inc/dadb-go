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

	fmt.Printf("%v\n", connectionResponse)

	err = WriteOpen(conn, 1, "shell:echo hello")
	if err != nil {
		panic(err)
	}

	packet, err := ReadPacket(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println(packet)

	conn.Close()
}
