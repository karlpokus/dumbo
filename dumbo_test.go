package dumbo

import (
	"bytes"
	"fmt"
	"testing"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%s", err)
	}
}

// a should match b
func TestPersistance(t *testing.T) {
  a := "this is nice"
	store := bytes.NewBuffer([]byte(a))
	data, err := New(store)
	fatal(t, err)
	err = data.Save(bytes.NewBuffer(data.gz))
	fatal(t, err)
  b := store.String()
  if a != b {
    fatal(t, fmt.Errorf("Expected %s and %s to be equal", a, b))
  }
}
