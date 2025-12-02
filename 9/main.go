package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run 9/main.go <string_to_unpack>")
		os.Exit(1)
	}

	toUnpack := os.Args[1]
	unpacked, err := Unpack(toUnpack)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(unpacked)
}

func Unpack(str string) (string, error) {
	toUnpack := []rune(str)
	unpacked := ""

	var symbol rune
	var times string
	unpackFlag := false
	escapeFlag := false

	for i, r := range toUnpack {
		if _, err := strconv.Atoi(string(r)); err == nil && !escapeFlag {
			times += string(r)
			unpackFlag = true

			if i+1 == len(toUnpack) {
				if e := unpackSymbol(symbol, &times, &unpacked, &unpackFlag); e != nil {
					return "", e
				}
			}
		} else {
			if r == rune('\\') {
				escapeFlag = true
			} else {
				if unpackFlag {
					if e := unpackSymbol(symbol, &times, &unpacked, &unpackFlag); e != nil {
						return "", e
					}
				}

				escapeFlag = false
				symbol = r
				unpacked += string(symbol)
			}
		}
	}

	return unpacked, nil
}

func unpackSymbol(symbol rune, times, unpacked *string, unpackFlag *bool) error {
	if symbol == rune(0) {
		return errors.New("invalid string to unpack")
	}

	n, _ := strconv.Atoi(*times)
	for k := 0; k < n-1; k++ {
		*unpacked += string(symbol)
	}

	*times = ""
	*unpackFlag = false

	return nil
}
