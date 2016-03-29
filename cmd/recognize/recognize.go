package main

import (
	"flag"
	"log"

	"github.com/armhold/gocarina"
	"os"
	"fmt"
)

var (
	boardFile string
	networkFile string
)

func init() {
	flag.StringVar(&boardFile, "board", "", "the letterpress board to read")
	flag.StringVar(&networkFile, "network", "", "the trained network file to use")
	flag.Parse()
}

func main() {
	if networkFile == "" || boardFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	board := gocarina.ReadUnknownBoard(boardFile)

	log.Printf("loading network...")
	network, err := gocarina.RestoreNetwork(networkFile)
	if err != nil {
		log.Fatal(err)
	}

	line := ""
	for i, tile := range board.Tiles {
		c := network.Recognize(tile.Reduced)
		line = line + fmt.Sprintf(" %c", c)

		// print them out shaped like a 5x5 letterpress board
		if (i + 1) % 5 == 0 {
			log.Printf(line)
			line = ""
		}
	}
}

