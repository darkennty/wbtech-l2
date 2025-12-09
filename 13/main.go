package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	dFlag *string
	fFlag *string
	sFlag *bool
)

// use: go run main.go -d ' ' -f 1,3-5 -s
func main() {
	dFlag = flag.String("d", "\t", "use given delimiter instead of '\\t' for field delimiter")
	fFlag = flag.String("f", "", "select only these fields; also print any line that contains no delimiter character, unless the -s option is specified")
	sFlag = flag.Bool("s", false, "do not print lines not containing delimiters")

	flag.Parse()

	fmt.Print("UNIX \"Cut\" command. Write a number N of lines to cut: ")
	in := bufio.NewReader(os.Stdin)

	var N int
	_, err := fmt.Scanf("%d\n", &N)
	if err != nil {
		log.Fatal("Error reading number of lines: ", err)
		return
	}

	fmt.Printf("Write %d lines to cut:\n", N)

	buffer := make([]string, 0)

	maxFieldsCount := 0

	for i := 0; i < N; i++ {
		command, _ := in.ReadString('\n')
		command = command[:len(command)-1]
		newFieldsCount := utf8.RuneCountInString(command)

		if newFieldsCount > maxFieldsCount {
			maxFieldsCount = newFieldsCount
		}

		buffer = append(buffer, command)
	}

	fields, err := getFields(*fFlag, maxFieldsCount)
	if err != nil {
		log.Fatal("Error parsing fields:", err)
		return
	}

	fmt.Println()
	fmt.Println("Result:")

	for _, line := range buffer {
		columns := strings.Split(line, *dFlag)
		if len(columns) == 1 && *sFlag {
			continue
		}

		if len(fields) == 0 {
			fmt.Println(line)
			continue
		}

		var outputFields []string
		for i := 1; i <= len(columns); i++ {
			if _, ok := fields[i]; ok || *fFlag == "" {
				outputFields = append(outputFields, columns[i-1])
			}
		}

		fmt.Println(strings.Join(outputFields, *dFlag))
	}
}

func getFields(f string, max int) (map[int]struct{}, error) {
	fields := make(map[int]struct{})
	if f == "" {
		return fields, nil
	}

	ranges := strings.Split(f, ",")
	for _, s := range ranges {
		n, err := strconv.Atoi(s)
		if err == nil {
			fields[n] = struct{}{}
		} else {
			bounds := strings.Split(s, "-")
			if len(bounds) == 2 {
				start, err1 := strconv.Atoi(bounds[0])
				end, err2 := strconv.Atoi(bounds[1])

				if err1 == nil && err2 == nil && start <= end {
					for i := start; i <= end; i++ {
						fields[i] = struct{}{}
					}
				} else if err1 == nil && err2 != nil {
					end = max
					for i := start; i <= end; i++ {
						fields[i] = struct{}{}
					}
				} else {
					return fields, errors.New("invalid field format")
				}
			} else {
				return fields, errors.New("invalid field format")
			}
		}
	}

	return fields, nil
}
