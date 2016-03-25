package gocarina

import (
	"fmt"
	_ "log"
	"math/rand"
	"time"
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Network struct {
	NumInputs     int         // total of bits in the image
	NumOutputs    int
	HiddenCount   int
	InputValues   []int       // image bits
	InputWeights  [][]float64 // weights from inputs -> hidden nodes
	HiddenOutputs []float64   // after feed-forward, what the hidden nodes output
	OutputWeights [][]float64 // weights from hidden nodes -> output nodes
	OutputValues  []float64   // after feed-forward, what the output nodes output
	OutputErrors  []float64
	HiddenErrors  []float64
}

func NewNetwork(charWidth, charHeight int) *Network {
	return &Network{NumInputs: charWidth * charHeight, HiddenCount: 25, NumOutputs: 8}
}

func (n *Network) assignRandomWeights() {

	// input -> hidden weights
	//
	for i := 0; i < n.NumInputs; i++ {
		weights := make([]float64, n.HiddenCount)

		for j := 0; j < len(weights); j++ {

			// we want the overall sum of weights to be < 1
			weights[j] = rand.Float64() / float64(n.NumInputs * n.HiddenCount)
		}

		n.InputWeights = append(n.InputWeights, weights)
	}


	// hidden -> output weights
	//
	for i := 0; i < n.HiddenCount; i++ {
		weights := make([]float64, n.NumOutputs)

		for j := 0; j < len(weights); j++ {

			// we want the overall sum of weights to be < 1
			weights[j] = rand.Float64() / float64(n.HiddenCount * n.NumOutputs)
		}

		n.OutputWeights = append(n.OutputWeights, weights)
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

func (n *Network) Save(filePath string) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(n)
	if err != nil {
		log.Fatalf("error encoding network: %s", err)
	}

	err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("error writing network to file: %s", err)
	}
}


func RestoreNetwork(filePath string) *Network {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading network file: %s", err)
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(b))

	var result Network
	err = decoder.Decode(&result)
	if err != nil {
		log.Fatalf("error decoding network: %s", err)
	}

	return &result
}
