package gocarina

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func WordsFrom(chars string) []string {
	chars = strings.ToLower(chars)

	file, err := os.Open("words-en.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var result []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if CanMakeWordFrom(word, chars) {
			result = append(result, word)
		}

		//fmt.Printf("word: %q\n", word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func CanMakeWordFrom(word string, chars string) bool {
	// turn char string into a pool of chars
	pool := []rune(chars)

	// iterate every char in word, and take them one at a time from pool
	for _, c := range word {
		var ok bool

		if pool, ok = takeOne(pool, c); !ok {
			// couldn't find c in p
			return false
		}
	}

	// found every letter
	return true
}

// if the char appears in the source string, remove it from source and return true
func takeOne(source []rune, char rune) ([]rune, bool) {
	for i, c := range source {
		if c == char {
			source = append(source[:i], source[i+1:]...)
			return source, true
		}
	}

	return source, false
}
