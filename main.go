package main

import (
	"golang.org/x/sync/errgroup"
	"net"
	"time"
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

	stream, err := connection.Open("shell:")
	if err != nil {
		panic(err)
	}

	eg := &errgroup.Group{}

	eg.Go(func() error {
		_, err = stream.Write([]byte("echo hello\n"))
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)

		return nil
	})

	eg.Go(func() error {
		packet, err := stream.readPacket()
		if err != nil {
			return err
		}
		print(string(packet.Payload))

		return nil
	})

	err = eg.Wait()
	if err != nil {
		panic(err)
	}
}
