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
	n.calculateOutputErrors('A')
	n.calculateFinalOutputs()
	n.calculateHiddenErrors()
	n.adjustOutputWeights()
	n.adjustInputWeights()
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

func TestCharToBinaryString(t *testing.T) {
	n := NewNetwork(25, 25)
	n.NumOutputs = 8

	cases := []struct {
		rune
		expected string
	}{
		{'0', "00110000"}, // we expect all strings to be left-padded with zeroes up to fill a width of n.NumOutputs
		{'a', "01100001"},
		{'A', "01000001"},
		{'z', "01111010"},
		{'Z', "01011010"},
	}

	for _, c := range cases {
		actual := n.charToBinaryString(c.rune)
		if actual != c.expected {
			t.Fatalf("for input %d, expected: %s, got: %s", c.rune, c.expected, actual)
		}
	}
}

func TestBinaryStringToInt(t *testing.T) {
	n := NewNetwork(25, 25)
	cases := []struct {
		s        string
		expected int64
	}{
		{"00110000", 48},
		{"01100001", 97},
		{"01000001", 65},
		{"01111010", 122},
		{"01011010", 90},
	}

	for _, c := range cases {
		actual := n.binaryStringToInt(c.s)
		if actual != c.expected {
			t.Fatalf("for input %s, expected: %d, got: %d", c.s, c.expected, actual)
		}
	}
}
