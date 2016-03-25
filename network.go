package gocarina

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Network struct {
	NumInputs     int // total of bits in the image
	NumOutputs    int
	HiddenCount   int
	InputValues   []int64     // image bits
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

// feed the image into the network
func (n *Network) assignInputs(img image.Image) {
	numPixels := img.Bounds().Dx() * img.Bounds().Dy()

	for i := 0; i < numPixels; i++ {
		col := pixelNumberToCol(i, img)
		row := pixelNumberToRow(i, img)
		pixel := pixelToBit(img.At(col, row))
		n.InputValues[i] = pixel
	}
}

func pixelNumberToCol(n int, img image.Image) int {
	return n % img.Bounds().Dx()
}

func pixelNumberToRow(n int, img image.Image) int {
	return n / img.Bounds().Dx()
}

func pixelToBit(c color.Color) int64 {
	if IsBlack(c) {
		return 0
	}

	return 1
}

func (n *Network) assignRandomWeights() {
	// input -> hidden weights
	//
	for i := 0; i < n.NumInputs; i++ {
		weights := make([]float64, n.HiddenCount)

		for j := 0; j < len(weights); j++ {

			// we want the overall sum of weights to be < 1
			weights[j] = rand.Float64() / float64(n.NumInputs*n.HiddenCount)
		}

		n.InputWeights = append(n.InputWeights, weights)
	}

	// hidden -> output weights
	//
	for i := 0; i < n.HiddenCount; i++ {
		weights := make([]float64, n.NumOutputs)

		for j := 0; j < len(weights); j++ {

			// we want the overall sum of weights to be < 1
			weights[j] = rand.Float64() / float64(n.HiddenCount*n.NumOutputs)
		}

		n.OutputWeights = append(n.OutputWeights, weights)
	}

	n.InputValues = make([]int64, n.NumInputs)
	n.OutputValues = make([]float64, n.NumOutputs)
	n.OutputErrors = make([]float64, n.NumOutputs)
	n.HiddenOutputs = make([]float64, n.NumOutputs)
	n.HiddenErrors = make([]float64, n.HiddenCount)
}

func (n *Network) calculateOutputErrors(trainedChar rune) {
	expected := float64(trainedChar)

	for i := 0; i < n.NumOutputs; i++ {
		n.OutputErrors[i] = (expected - n.OutputValues[i]) * (1.0 - n.OutputValues[i]) * n.OutputValues[i]
	}
}

func (n *Network) calculateHiddenErrors() {
	for i := 0; i < len(n.HiddenOutputs); i++ {
		sum := float64(0)

		for j := 0; j < len(n.OutputErrors); j++ {
			sum += n.OutputErrors[j] * n.OutputWeights[i][j]
		}

		n.HiddenErrors[i] = n.HiddenOutputs[i] * (1 - n.HiddenOutputs[i]) * sum
	}
}

func (n *Network) adjustOutputWeights() {
	for i := 0; i < len(n.HiddenOutputs); i++ {
		for j := 0; j < n.NumOutputs; j++ {
			n.OutputWeights[i][j] += n.OutputErrors[j] * n.HiddenOutputs[i]
		}
	}
}

func (n *Network) adjustInputWeights() {
	for i := 0; i < n.NumInputs; i++ {
		for j := 0; j < n.HiddenCount; j++ {
			//fmt.Printf("i: %d, j: %d, len(n.InputWeights): %d, len(n.HiddenErrors): %d, len(n.InputValues): %d\n", i, j, len(n.InputWeights), len(n.HiddenErrors), len(n.InputValues))
			n.InputWeights[i][j] += n.HiddenErrors[j] * float64(n.InputValues[i])
		}
	}
}

func (n *Network) calculateHiddenOutputs() {
	for i := 0; i < len(n.HiddenOutputs); i++ {
		sum := float64(0)

		for j := 0; j < len(n.InputValues); j++ {
			sum += float64(n.InputValues[j]) * n.InputWeights[j][i]
		}

		n.HiddenOutputs[i] = sigmoid(sum)
	}
}

func (n *Network) calculateFinalOutputs() {
	for i := 0; i < n.NumOutputs; i++ {
		sum := float64(0)

		for j := 0; j < len(n.HiddenOutputs); j++ {
			sum += n.HiddenOutputs[j] * n.OutputWeights[j][i]
		}

		n.OutputValues[i] = sigmoid(sum)
	}
}

// function that maps its input to a range between 0..1
// mathematically it's supposed to be asymptotic, but large values of x may round up to 1
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
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

func (n *Network) intToBinaryString(i int64) string {
	// we want to pad with n.NumOutputs number of zeroes, so create a dynamic format for Sprintf
	format := fmt.Sprintf("%%0%db", n.NumOutputs)
	return fmt.Sprintf(format, i)
}

func (n *Network) charToBinaryString(c rune) string {
	return n.intToBinaryString(int64(c))
}

func (n *Network) binaryStringToInt(s string) int64 {
	result, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		log.Fatalf("error converting binary string %s to int: %s", s, err)
	}

	return result
}
