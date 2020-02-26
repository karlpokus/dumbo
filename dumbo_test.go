package dumbo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%s", err)
	}
}

// a should match b
func TestPersistance(t *testing.T) {
	fpath := "testdata/nice.json"
	a, err := ioutil.ReadFile(fpath)
	fatal(t, err)

	data, err := New(fpath)
	fatal(t, err)
	err = data.Save(bytes.NewBuffer(data.gz))
	fatal(t, err)

	b, err := ioutil.ReadFile(fpath)
	fatal(t, err)
	if !bytes.Equal(a, b) {
		fatal(t, fmt.Errorf("Expected %s and %s to be equal", a, b))
	}
}
