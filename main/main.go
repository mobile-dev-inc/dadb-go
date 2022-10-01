package main

import (
	"dadb/adbd"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5555")
	if err != nil {
		panic(err)
	}

	dadb, err := adbd.CreateDadb(conn)
	if err != nil {
		panic(err)
	}

	stream, err := dadb.Open("shell:")
	if err != nil {
		panic(err)
	}

	eg := &errgroup.Group{}

	eg.Go(func() error {
		b := make([]byte, 1)
		for {
			_, err := os.Stdin.Read(b)
			if err != nil {
				return err
			}
			_, err = stream.Write(b)
			if err != nil {
				return err
			}
		}
	})

	eg.Go(func() error {
		for {
			buffer := make([]byte, 1024)
			n, err := stream.Read(buffer)
			if err != nil {
				return err
			}
			print(string(buffer[:n]))
		}
	})

	err = eg.Wait()
	if err != nil {
		panic(err)
	}
}
