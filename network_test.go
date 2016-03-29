package gocarina

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNetwork(t *testing.T) {
	n := NewNetwork(25 * 25)
	n.calculateHiddenOutputs()
	n.calculateOutputErrors('A')
	n.calculateFinalOutputs()
	n.calculateHiddenErrors()
	n.adjustOutputWeights()
	n.adjustInputWeights()
}

func TestSaveRestore(t *testing.T) {
	n := NewNetwork(25 * 25)
	n.assignRandomWeights()

	f, err := ioutil.TempFile("", "network")
	if err != nil {
		t.Fatal(err)
	}

	n.Save(f.Name())
	restored, err := RestoreNetwork(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(n, restored) {
		t.Fatalf("expected: %+v, got %+v", n, restored)
	}
}

func TestSigmoid(t *testing.T) {
	xvals := []float64{-100000.0, -10000.0, -1000.0, -100.0, -10.0, 0.0, 0.1, 0.01, 0.001, 1.0, 10.0, 100.0, 1000.0, 10000.0, 100000.0}

	for _, x := range xvals {
		y := sigmoid(x)

		if y < 0.0 || y > 1.0 {
			t.Fatalf("for input %f got %f, should be in (0..1.0)", x, y)
		}
	}
}

func TestRuneToArrayOfInts(t *testing.T) {
	n := NewNetwork(25 * 25)

	expected := []int{0, 1, 0, 0, 0, 0, 0, 1}
	actual := n.runeToArrayOfInts('A')

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %+v, got: %+v", expected, actual)
	}
}
