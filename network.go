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

func NewNetwork(numInputs int) *Network {
	// NB: NumOutputs effectively constrains the range of chars that are recognizable
	n := &Network{NumInputs: numInputs, HiddenCount: 250, NumOutputs: 8}

	n.InputValues = make([]int64, n.NumInputs)
	n.OutputValues = make([]float64, n.NumOutputs)
	n.OutputErrors = make([]float64, n.NumOutputs)
	n.HiddenOutputs = make([]float64, n.NumOutputs)
	n.HiddenErrors = make([]float64, n.HiddenCount)

	n.assignRandomWeights()

	return n
}

// Attempt to recognize the character displayed on the given image.
func (n *Network) Recognize(img image.Image) rune {
	img = BlackWhiteImage(img)

	n.assignInputs(img)
	n.calculateHiddenOutputs()
	n.calculateFinalOutputs()

	// quantize output values
	bitstring := ""
	for _, v := range n.OutputValues {
		//log.Printf("v: %f", v)
		bitstring += strconv.Itoa(Round(v))
	}

	asciiCode, err := strconv.ParseInt(bitstring, 2, 16)
	if err != nil {
		log.Fatalf("error in ParseInt for %s: ", err)
	}

	//log.Printf("returning bitstring: %s", bitstring)
	return rune(asciiCode)
}

// can't believe this isn't in the stdlib!
func Round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

// feed the image into the network
func (n *Network) assignInputs(img image.Image) {
	numPixels := img.Bounds().Dx() * img.Bounds().Dy()
	if numPixels != n.NumInputs {
		log.Fatalf("expected %d inputs, got %d", n.NumInputs, numPixels)
	}
	//log.Printf("numPixels: %d", numPixels)

	i := 0
	for row := img.Bounds().Min.Y; row < img.Bounds().Max.Y; row++ {
		for col := img.Bounds().Min.X; col < img.Bounds().Max.X; col++ {
			pixel := pixelToBit(img.At(col, row))
			n.InputValues[i] = pixel
			i++
		}
	}

	if i != n.NumInputs {
		log.Fatalf("expected i to be: %d, was: %d", n.NumInputs, i)
	}
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
}

func (n *Network) calculateOutputErrors(r rune) {
	accumError := 0.0
	arrayOfInts := n.runeToArrayOfInts(r)

	// NB: binaryString[i] will return bytes, not a rune. range does the right thing
	for i, digit := range arrayOfInts {
		//log.Printf("digit: %d", digit)

		digitAsFloat := float64(digit)
		err := (digitAsFloat - n.OutputValues[i]) * (1.0 - n.OutputValues[i]) * n.OutputValues[i]
		n.OutputErrors[i] = err
		accumError += err * err
		//log.Printf("accumError: %.10f", accumError)
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
			val := n.HiddenOutputs[j] * n.OutputWeights[j][i]
			sum += val
			//log.Printf("val: %f", val)
		}

		//log.Printf("sum: %f", sum)
		n.OutputValues[i] = sigmoid(sum)
	}
}

// function that maps its input to a range between 0..1
// mathematically it's supposed to be asymptotic, but large values of x may round up to 1
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func (n *Network) Train(img image.Image, r rune) {
	// feed the image data forward through the network to obtain a result
	//
	n.assignInputs(img)
	n.calculateHiddenOutputs()
	n.calculateFinalOutputs()

	// propagate the error correction backward through the net
	//
	n.calculateOutputErrors(r)
	n.calculateHiddenErrors()
	n.adjustOutputWeights()
	n.adjustInputWeights()
}

func (n *Network) Save(filePath string) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(n)
	if err != nil {
		return fmt.Errorf("error encoding network: %s", err)
	}

	err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing network to file: %s", err)
	}

	return nil
}

func RestoreNetwork(filePath string) (*Network, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading network file: %s", err)
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(b))

	var result Network
	err = decoder.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding network: %s", err)
	}

	return &result, nil
}

// map a rune char to an array of int, representing its unicode codepoint in binary
// 'A' => 65 => []int {0, 1, 0, 0, 0, 0, 0, 1}
// result is zero-padded to n.NumOutputs
//
func (n *Network) runeToArrayOfInts(r rune) []int {
	var result []int = make([]int, n.NumOutputs)

	codePoint := int64(r) // e.g. 65

	// we want to pad with n.NumOutputs number of zeroes, so create a dynamic format for Sprintf
	format := fmt.Sprintf("%%0%db", n.NumOutputs)
	binaryString := fmt.Sprintf(format, codePoint) // e.g. "01000001"

	// must use range: array indexing of strings returns bytes
	for i, v := range binaryString {
		if v == '0' {
			result[i] = 0
		} else {
			result[i] = 1
		}
	}
	return result
}
