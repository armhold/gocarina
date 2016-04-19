package main

import (
	"flag"
	"log"

	"github.com/armhold/gocarina"
)

var (
	fromFile string
	toFile   string
	maxIter  int
)

func init() {
	flag.StringVar(&fromFile, "load", "", "to load network from a saved file")
	flag.StringVar(&toFile, "save", "ocr.save", "to save network to a file")
	flag.IntVar(&maxIter, "max", 500, "max number of training iterations")

	flag.Parse()
}

func main() {
	log.SetFlags(0)

	// do this first, so we have tile boundaries to create the network
	m := gocarina.ReadKnownBoards()

	exampleTile := m['A'].Reduced
	pixelCount := exampleTile.Bounds().Dx() * exampleTile.Bounds().Dy()
	numInputs := pixelCount

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

	// save files for debugging
	for _, tile := range m {
		tile.SaveReducedTile()
	}

	for i := 0; i < maxIter; i++ {
		for r, tile := range m {
			//log.Printf("training: %c\n", r)
			n.Train(tile.Reduced, r)
		}

		if allCorrect(m, n) {
			log.Printf("success took %d iterations", i+1)
			break
		}
	}

	if toFile != "" {
		n.Save(toFile)
	}

	// show details on success/failure
	count := 0
	correct := 0
	for r, tile := range m {
		recognized := n.Recognize(tile.Reduced)
		count++

		if recognized == r {
			correct++
		} else {
			log.Printf("failure: tile recognized as: %c, should be: %c", recognized, r)
		}
	}

	successPercent := float64(correct) / float64(count) * 100.0
	log.Printf("success rate: %d/%d => %%%.2f", correct, count, successPercent)
}

// fail-fast way to detect if network has 100% success rate yet
func allCorrect(m map[rune]*gocarina.Tile, n *gocarina.Network) bool {
	for r, tile := range m {
		recognized := n.Recognize(tile.Reduced)
		if recognized != r {
			return false
		}
	}

	return true
}
