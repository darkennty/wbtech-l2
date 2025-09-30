package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	//fmt.Printf("Value: %#v, type: %T\n", err, err) // Value: (*fs.PathError)(nil), type: *fs.PathError
	fmt.Println(err == nil)
}
