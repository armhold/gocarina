package gocarina

import (
	"testing"
	"io/ioutil"
	"reflect"
)

func TestNetwork(t *testing.T) {

	n := NewNetwork(25, 25)
	n.assignRandomWeights()

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
