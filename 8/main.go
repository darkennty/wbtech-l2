package main

import (
	"fmt"
	ntp "github.com/darkennty/ntp-time"
	"os"
)

func main() {
	t, err := ntp.GetTime()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Current time: %s", t)
}
