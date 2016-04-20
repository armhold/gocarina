package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/armhold/gocarina"
)

var (
	boardFile   string
	networkFile string
	showWords   bool
)

func init() {
	var usage = func() {
		u := `Usage: recognize [-w] game_board.png

Supply a game board file as an argument (e.g. board-images/board1.png),
and 'recognize' will print the letters it contains.

If -w is specified, the set of words that can be formed by the game board letters
is also printed.

`
		fmt.Fprintf(os.Stderr, u)
	}

	flag.Usage = usage

	flag.StringVar(&networkFile, "network", "ocr.save", "the trained network file to use")
	flag.BoolVar(&showWords, "w", false, "show list of words that can be made from the given board")
	flag.Parse()

	// remaining args after flags are parsed
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	boardFile = args[0]
}

func main() {
	log.SetFlags(0)

	if networkFile == "" || boardFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	board := gocarina.ReadUnknownBoard(boardFile)

	//log.Printf("loading network...")
	network, err := gocarina.RestoreNetwork(networkFile)
	if err != nil {
		log.Fatal(err)
	}

	line := ""
	allLetters := ""
	for i, tile := range board.Tiles {
		c := network.Recognize(tile.Reduced)
		line = line + fmt.Sprintf(" %c", c)
		allLetters = allLetters + string(c)

		// print them out shaped like a 5x5 letterpress board
		if (i+1)%5 == 0 {
			log.Printf(line)
			line = ""
		}
	}

	if showWords {
		log.Printf("\n\n")

		words := gocarina.WordsFrom(allLetters)
		for _, word := range words {
			log.Println(word)
		}
	}
}
