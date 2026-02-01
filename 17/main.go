package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var timeout *int
	timeout = pflag.Int("timeout", 10, "Sets a connection timeout (in seconds).")
	pflag.Parse()

	if len(pflag.Args()) < 2 {
		fmt.Println("Usage: go run main.go <URL> <PORT> [--timeout=<DURATION>]")
		return
	}

	URL := os.Args[1]
	PORT := os.Args[2]

	fmt.Printf("Trying %s...\n", URL)

	conn, err := net.DialTimeout("tcp", URL+":"+PORT, time.Duration(*timeout)*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Connected to %s.\nEscape character is '^D'.\n", URL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalHandler(cancel)

	go read(conn, ctx, cancel)
	go write(conn, cancel)

	<-ctx.Done()

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func read(conn net.Conn, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	if _, err := io.Copy(os.Stdout, conn); err != nil {
		select {
		case <-ctx.Done():
		default:
			log.Println(err)
		}
	} else {
		log.Println("Connection closed by remote host.")
	}
}

func write(conn net.Conn, cancel context.CancelFunc) {
	defer cancel()

	if _, err := io.Copy(conn, os.Stdin); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}
}

func signalHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}
