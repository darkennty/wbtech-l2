package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Local server created to test main.go

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024) // temp
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v\n", err)
			}
			break
		}

		fmt.Printf("Received %d bytes: %s", n, string(buf[:n]))

		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("Write error: %v\n", err)
			break
		}
	}
	fmt.Println("Connection closed.")
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Server listening on :8080")

	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection accepted.")

	handleConnection(conn)
}
