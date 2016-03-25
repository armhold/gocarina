package gocarina

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNetwork(t *testing.T) {

	n := NewNetwork(25, 25)
	n.assignRandomWeights()
	n.calculateHiddenOutputs()
	n.calculateOutputErrors()

	n.printInputWeights()
}

func TestSaveRestore(t *testing.T) {
	n := NewNetwork(25, 25)
	n.assignRandomWeights()

	f, err := ioutil.TempFile("", "network")
	if err != nil {
		t.Fatal(err)
	}

	n.Save(f.Name())
	restored := RestoreNetwork(f.Name())

	if !reflect.DeepEqual(n, restored) {
		t.Errorf("expected: %+v, got %+v", n, restored)
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
