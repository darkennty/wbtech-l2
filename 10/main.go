package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	kFlag *int
	nFlag *bool
	rFlag *bool
	uFlag *bool
	bFlag *bool
	cFlag *bool
	hFlag *bool
	MFlag *bool
	tabs  = regexp.MustCompile(`\s+`)
)

func expandShortBoolFlags(args []string, boolFlags map[rune]struct{}) []string {
	if len(args) == 0 {
		return args
	}
	out := []string{args[0]}
	for i := 1; i < len(args); i++ {
		a := args[i]

		if a[0] == '-' && len(a) > 1 && a[1] != '-' {
			name := a[1:]
			ok := true
			var parts []string

			for _, r := range name {
				if _, exists := boolFlags[r]; !exists {
					ok = false
					break
				}
				parts = append(parts, "-"+string(r))
			}

			if ok {
				out = append(out, parts...)
				continue
			}
		}

		out = append(out, a)
	}

	return out
}

func main() {
	boolShorts := map[rune]struct{}{
		'n': {}, 'r': {}, 'u': {}, 'b': {}, 'c': {}, 'h': {}, 'M': {},
	}
	os.Args = expandShortBoolFlags(os.Args, boolShorts)

	kFlag = flag.Int("k", 1, "sort by comparing k-th column")
	nFlag = flag.Bool("n", false, "numeric sort (compare according to string numerical value)")
	rFlag = flag.Bool("r", false, "reverse the result of comparisons")
	uFlag = flag.Bool("u", false, "output only unique lines")
	bFlag = flag.Bool("b", false, "remove trailing blanks")
	cFlag = flag.Bool("c", false, "check for sorted input; do not sort")
	hFlag = flag.Bool("h", false, "compare human readable numbers (e.g., 2K 1G)")
	MFlag = flag.Bool("M", false, "month sort (compare <unknown> < 'JAN' < ... < 'DEC')")

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Error: input file is not specified")
		return
	}

	if *hFlag {
		*nFlag = true
	}

	*kFlag--

	inputFileName := flag.Arg(0)

	if *cFlag {
		sorted, err := checkSorted(inputFileName)
		if err != nil {
			log.Fatalf("Error checking file: %v", err)
		}
		if sorted {
			fmt.Println("Input file is sorted")
		} else {
			fmt.Println("Input file is not sorted")
		}
		return
	}

	chunkFiles, err := createInitialRuns(inputFileName)
	if err != nil {
		log.Fatalf("Error creating initial runs: %v", err)
	}
	defer cleanup(chunkFiles)

	err = mergeRuns(chunkFiles, os.Stdout)
	if err != nil {
		log.Fatalf("Error merging runs: %v", err)
	}
}

func createInitialRuns(inputFileName string) ([]string, error) {
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	chunkFiles := make([]string, 0)
	chunk := make([]string, 0)
	chunkSize := 100000

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) >= chunkSize {
			if err = sortAndWriteChunk(&chunk, &chunkFiles); err != nil {
				return nil, err
			}
		}
	}

	if len(chunk) > 0 {
		if err = sortAndWriteChunk(&chunk, &chunkFiles); err != nil {
			return nil, err
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input file: %w", err)
	}

	return chunkFiles, nil
}

func sortAndWriteChunk(chunk *[]string, chunkFiles *[]string) error {
	sort.Slice(*chunk, func(i, j int) bool {
		return CompareStrings((*chunk)[i], (*chunk)[j])
	})

	tempFile, err := os.Create("sort-chunk.txt")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	writer := bufio.NewWriter(tempFile)
	for _, line := range *chunk {
		if _, err = fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}
	if err = writer.Flush(); err != nil {
		return err
	}

	*chunkFiles = append(*chunkFiles, tempFile.Name())
	*chunk = nil
	return nil
}

type HeapItem struct {
	line    string
	scanner *bufio.Scanner
}

type PriorityQueue []*HeapItem

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return CompareStrings(pq[i].line, pq[j].line)
}
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*HeapItem)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func mergeRuns(chunkFiles []string, output io.Writer) error {
	if len(chunkFiles) == 0 {
		return nil
	}

	var readers []*os.File
	defer func() {
		for _, r := range readers {
			r.Close()
		}
	}()

	pq := make(PriorityQueue, 0, len(chunkFiles))

	for _, fileName := range chunkFiles {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		readers = append(readers, file)

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			item := &HeapItem{
				line:    scanner.Text(),
				scanner: scanner,
			}
			pq = append(pq, item)
		}
	}

	heap.Init(&pq)

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	var lastWritten string
	isFirst := true
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*HeapItem)

		if *uFlag {
			if isFirst || item.line != lastWritten {
				fmt.Fprintln(writer, item.line)
				lastWritten = item.line
				isFirst = false
			}
		} else {
			fmt.Fprintln(writer, item.line)
		}

		if item.scanner.Scan() {
			item.line = item.scanner.Text()
			heap.Push(&pq, item)
		} else if err := item.scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}

func checkSorted(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var prevLine string
	isFirst := true

	for scanner.Scan() {
		currentLine := scanner.Text()
		if isFirst {
			prevLine = currentLine
			isFirst = false
			continue
		}

		if !CompareStrings(prevLine, currentLine) {
			return false, nil
		}
		prevLine = currentLine
	}

	return true, scanner.Err()
}

func cleanup(files []string) {
	for _, file := range files {
		os.Remove(file)
	}
}

func CompareStrings(strA, strB string) bool {
	if *bFlag {
		strA = strings.TrimRight(strA, " \t")
		strB = strings.TrimRight(strB, " \t")
	}

	first := tabs.Split(strA, -1)
	second := tabs.Split(strB, -1)

	months := map[string]int{
		"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4,
		"MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8,
		"SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
		"UNKNOWN": 0,
	}

	k := *kFlag
	if k < 0 {
		k = 0
	}

	var firstKey, secondKey string
	if len(first) > k {
		firstKey = first[k]
	}
	if len(second) > k {
		secondKey = second[k]
	}

	if *MFlag {
		fMon, ok1 := months[strings.ToUpper(firstKey)]
		sMon, ok2 := months[strings.ToUpper(secondKey)]

		if !ok1 {
			fMon = months["UNKNOWN"]
		}
		if !ok2 {
			sMon = months["UNKNOWN"]
		}

		if fMon != sMon {
			if *rFlag {
				return fMon > sMon
			}
			return fMon < sMon
		}
	}

	if *nFlag {
		var aVal, bVal int64

		if *hFlag {
			a, err1 := humanReadableToBytes(firstKey)
			b, err2 := humanReadableToBytes(secondKey)
			if err1 == nil && err2 == nil {
				aVal, bVal = a, b
			}
		} else {
			a, err1 := strconv.ParseInt(firstKey, 10, 64)
			b, err2 := strconv.ParseInt(secondKey, 10, 64)
			if err1 == nil && err2 == nil {
				aVal, bVal = a, b
			}
		}

		if aVal != bVal {
			if *rFlag {
				return aVal > bVal
			}
			return aVal < bVal
		}
	}

	if strA != strB {
		if *rFlag {
			return strA > strB
		}
		return strA < strB
	}

	return true
}

func humanReadableToBytes(s string) (int64, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return 0, nil
	}

	var multiplier int64 = 1
	suffix := s[len(s)-1]

	if suffix >= 'A' && suffix <= 'Z' {
		switch suffix {
		case 'K':
			multiplier = 1024
		case 'M':
			multiplier = 1024 * 1024
		case 'G':
			multiplier = 1024 * 1024 * 1024
		case 'T':
			multiplier = 1024 * 1024 * 1024 * 1024
		default:
			multiplier = 1
		}
		if multiplier > 1 {
			s = s[:len(s)-1]
		}
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid human readable number format: %w", err)
	}
	return val * multiplier, nil
}
