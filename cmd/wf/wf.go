package main

import (
	"fmt"
	"github.com/armhold/gocarina"
	"log"
	"os"
)

var (
	allLetters string
)

func init() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: wf [letters], e.g.: wf ttpwnredkocnutlsrcowodwua")
	}

	allLetters = os.Args[1]
}

// wf- "words from": prints list of words that can be formed from a given list of characters.
// Useful for deriving words manually, without bothering with game board images.
func main() {
	log.SetFlags(0)
	fmt.Printf("words from %s\n", allLetters)

	words := gocarina.WordsFrom(allLetters)
	for _, word := range words {
		fmt.Println(word)
	}
}
