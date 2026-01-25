package main

import (
	"fmt"
	"os"
	"strconv"
	"wbtech_l2/16/loader"
	"wbtech_l2/16/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <URL> [RECURSION_DEPTH]")
		return
	}

	url := os.Args[1]
	var queue []*parser.Node

	var N int
	var err error
	if len(os.Args) > 2 {
		depth := os.Args[2]
		N, err = strconv.Atoi(depth)
		if err != nil {
			N = -1
		} else {
			N += 1 // to load current page + pages on depth of n
		}
	} else {
		N = -1
	}

	queue = append(queue, &parser.Node{Url: url, N: N})
	visited := make(map[string]struct{})

	for len(queue) > 0 {
		loader.Load(&queue, visited)
		queue = queue[1:]
	}
}
