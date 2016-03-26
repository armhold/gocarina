package main

import (
	"flag"
	"log"

	"github.com/armhold/gocarina"
)

var (
	fromFile string
	toFile string
	iter int
)

func init() {
	flag.StringVar(&fromFile, "load", "", "to load network from a saved file")
	flag.StringVar(&toFile, "save", "", "to save network to a file")
	flag.IntVar(&iter, "iter", 100, "number of training iterations")

	flag.Parse()
}

func main() {
	// do this first, so we have tile boundaries to create the network
	m := gocarina.ProcessGameBoards()

	targetRune := 'A'
	tile := m[targetRune]
	numInputs := tile.Bounds().Dx() * tile.Bounds().Dy()

	var n *gocarina.Network
	var err error

	if fromFile != "" {
		log.Printf("loading network...")
		n, err = gocarina.RestoreNetwork(fromFile)
		if err != nil {
			log.Fatal(err)
		}

		if n.NumInputs != numInputs {
			log.Fatalf("loaded network has %d inputs, tile has %d", n.NumInputs, numInputs)
		}


	} else {
		log.Printf("creating new network...")
		n = gocarina.NewNetwork(numInputs)
	}


	for i := 0; i < iter; i++ {
		for r, tile := range m {
			log.Printf("training: %c\n", r)
			n.Train(tile, r)
		}
	}

	if toFile != "" {
		n.Save(toFile)
	}


	result := n.Recognize(tile)
	success := result == targetRune
	log.Printf("tile recognized as: %c, success: %t", result, success)
}
