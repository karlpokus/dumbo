package dumbo

import (
	"compress/gzip"
	"crypto/sha1"
	"io"
)

// compress compresses r to w
func compress(w io.Writer, r io.Reader) error {
	gz := gzip.NewWriter(w)
	if _, err := io.Copy(gz, r); err != nil {
		return err
	}
	return gz.Close()
}

// Decompress decompresses r to w
func Decompress(w io.Writer, r io.Reader) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, gz); err != nil {
		return err
	}
	return gz.Close()
}

// hash returns b hashed
func hash(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return h.Sum(nil)
}
