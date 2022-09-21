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

	err, connectionResponse := connect(conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", connectionResponse)

	err = writeOpen(conn, 1, "shell:echo hello")
	if err != nil {
		panic(err)
	}

	packet, err := readPacket(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println(packet)

	packet, _ = readPacket(conn)
	fmt.Println(string(packet.Payload))

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
