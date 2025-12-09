package main

import (
	"fmt"
	"sort"
)

func main() {
	dict := []string{"тяпка", "пятак", "пятка", "столик", "листок", "слиток", "стол"}
	anagrams := getAnagrams(dict)

	for k, v := range anagrams {
		fmt.Println(k, "->", v)
	}
}

func getAnagrams(dict []string) map[string][]string {
	anagrams := make(map[string][]string, len(dict))
	fromSortedToKeyMap := make(map[string]string, len(dict))

	for _, word := range dict {
		sortedWord := sortString(word)
		key, exists := fromSortedToKeyMap[sortedWord]

		if !exists {
			key = word
			fromSortedToKeyMap[sortedWord] = key
		}

		anagrams[key] = append(anagrams[key], word)
	}

	for k, v := range anagrams {
		if len(v) < 2 {
			delete(anagrams, k)
		} else {
			sort.Slice(anagrams[k], func(i, j int) bool {
				return anagrams[k][i] < anagrams[k][j]
			})
		}
	}

	return anagrams
}

func sortString(s string) string {
	runes := []rune(s)

	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	return string(runes)
}
