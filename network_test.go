package gocarina

import (
	"testing"
)

func TestNetwork(t *testing.T) {

	n := NewNetwork(25, 25)
	n.assignRandomWeights()

	n.printInputWeights()
}
