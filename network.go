package gocarina

import (
	"fmt"
	_ "log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Network struct {
	NumInputs     int // total of bits in the image
	HiddenCount   int
	InputValues   []int       // image bits
	InputWeights  [][]float64 // weights from inputs -> hidden nodes
	HiddenOutputs []float64   // after feed-forward, what the hidden nodes output
	OutputWeights []float64   // weights from hidden nodes -> output nodes
	OutputValues  []float64   // after feed-forward, what the output nodes output
	OutputErrors  []float64
	HiddenErrors  []float64
}

func NewNetwork(charWidth, charHeight int) *Network {
	return &Network{NumInputs: charWidth * charHeight, HiddenCount: 25}
}

func (n *Network) assignRandomWeights() {

	for i := 0; i < n.NumInputs; i++ {
		weights := make([]float64, n.HiddenCount)

		for j := 0; j < len(weights); j++ {
			weights[j] = rand.Float64()
		}

		n.InputWeights = append(n.InputWeights, weights)
	}

}

func (n *Network) printInputWeights() {
	for _, weights := range n.InputWeights {

		for _, w := range weights {
			fmt.Printf("%f ", w)
		}

		fmt.Println()
	}
}
