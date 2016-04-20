package gocarina

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

// WordsFrom returns a slice of dictionary words that can be constructed from the given chars.
func WordsFrom(chars string) []string {
	chars = strings.ToLower(chars)

	file, err := os.Open("words-en.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var result []string

	// iterate every word in the file. NB: words are already lower-cased in the file.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if CanMakeWordFrom(word, chars) {
			result = append(result, word)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Sort(ByWordLength(result))

	return result
}

// CanMakeWordFrom returns true if the characters from 'chars' can be re-ordered to form 'word', else false.
// Leftover letters are OK, but individual letters cannot be re-used. If a given letter is needed multiple times
// (e.g. 'door' needs two o's), then the letter must appear multiple times in 'chars'.
func CanMakeWordFrom(word string, chars string) bool {
	pool := []rune(chars)

	// iterate every char in word, and take them one at a time from pool
	for _, c := range word {
		var ok bool

		if pool, ok = takeOne(pool, c); !ok {
			// couldn't find c in pool
			return false
		}
	}

	// found every letter
	return true
}

// takeOne will remove one instance of the given char from pool. It returns the (possibly modified) slice,
// and a boolean to indicate whether the char was found or not.
func takeOne(pool []rune, char rune) ([]rune, bool) {
	for i, c := range pool {
		if c == char {
			result := append(pool[:i], pool[i+1:]...)
			return result, true
		}
	}

	return pool, false
}

// Sort by word length descending, then alphabetical ascending. So bigger words come first,
// but equal-length words are sub-sorted alphabetically.
type ByWordLength []string

func (w ByWordLength) Len() int      { return len(w) }
func (w ByWordLength) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w ByWordLength) Less(i, j int) bool {
	ri := []rune(w[i])
	rj := []rune(w[j])

	// first sort on word-length, descending
	if len(ri) > len(rj) {
		return true
	}

	if len(ri) < len(rj) {
		return false
	}

	// lengths are equal, so sort secondarily by alphabet
	for i, _ := range ri {
		if ri[i] == rj[i] {
			continue
		}

		if ri[i] < rj[i] {
			return true
		}

		return false
	}

	// if we get here, all chars are equal, but "x" is still not < "x", so return false
	return false
}
