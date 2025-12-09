package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	AFlag *int
	BFlag *int
	CFlag *int
	cFlag *bool
	iFlag *bool
	vFlag *bool
	FFlag *bool
	nFlag *bool
)

// grep: grep [опции] шаблон [<путь к файлу или папке>]
// use: go run main.go -B 2 -A 4 -c -n -F -v -i 'auth' testdata.txt
func main() {
	AFlag = flag.Int("A", 0, "Print NUM lines of trailing context after matching lines.")
	BFlag = flag.Int("B", 0, "Print NUM lines of leading context before matching lines.")
	CFlag = flag.Int("C", 0, "Print NUM lines of output context. Equivalent to '-A NUM -B NUM'.")
	cFlag = flag.Bool("c", false, "Suppress normal output; instead print a count of matching lines for each input file. With the '-v' option, count non-matching lines.") // количество строк, которые соответствуют шаблону
	iFlag = flag.Bool("i", false, "Ignore case distinctions in patterns and input data, so that characters that differ only in case match each other.")
	vFlag = flag.Bool("v", false, "Invert the sense of matching, to select non-matching lines.")
	FFlag = flag.Bool("F", false, "Interpret patterns as fixed strings, not regular expressions.")
	nFlag = flag.Bool("n", false, "Prefix each line of output with the 1-based line number within its input file.")

	if len(os.Args) < 3 {
		log.Fatal("Error: pattern and input filenames are required")
		return
	}

	inputFileName := os.Args[len(os.Args)-1]
	pattern := os.Args[len(os.Args)-2]

	os.Args = os.Args[:len(os.Args)-2]

	flag.Parse()

	if *AFlag == 0 && *CFlag != 0 {
		*AFlag = *CFlag
	}

	if *BFlag == 0 && *CFlag != 0 {
		*BFlag = *CFlag
	}

	if len(os.Args) < flag.NFlag()+1 {
		log.Fatal("Error: pattern and input filenames are required")
		return
	}

	matchingLinesCounter := 0

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal("can't open the file: ", err)
	}
	defer inputFile.Close()

	strings.Trim(pattern, "'")

	if *iFlag {
		pattern = strings.ToLower(pattern)
	}

	if *FFlag {
		pattern = strings.Join([]string{"^", pattern, "$"}, "")
	}

	re := regexp.MustCompile(pattern)

	BSlice := make([]string, 0, *BFlag)
	toPrint := 0

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if *iFlag {
			line = strings.ToLower(line)
		}

		match := re.MatchString(line)

		if *vFlag {
			match = !match
		}

		if toPrint != 0 {
			if match {
				toPrint = *AFlag + 1
			}

			matchingLinesCounter++
			if !*cFlag {
				if *nFlag {
					fmt.Printf("%d. %s\n", matchingLinesCounter, line)
				} else {
					fmt.Println(line)
				}
			}

			toPrint--
			continue
		}

		if match {
			if *AFlag != 0 {
				toPrint = *AFlag
			}

			if *BFlag != 0 {
				for i := 0; i < len(BSlice); i++ {
					matchingLinesCounter++

					if !*cFlag {
						if *nFlag {
							fmt.Printf("%d. %s\n", matchingLinesCounter, BSlice[i])
						} else {
							fmt.Println(BSlice[i])
						}
					}
				}

				BSlice = make([]string, 0, *BFlag)
			}

			matchingLinesCounter++
			if !*cFlag {
				if *nFlag {
					fmt.Printf("%d. %s\n", matchingLinesCounter, line)
				} else {
					fmt.Println(line)
				}
			}
		} else {
			if *BFlag != 0 {
				if len(BSlice) == *BFlag {
					BSlice = BSlice[1:*BFlag:*BFlag]
				}
				BSlice = append(BSlice, line)
			}
		}
	}

	if *cFlag {
		fmt.Println(matchingLinesCounter)
	}
}
