package dumbo

import (
	"crypto/sha1"
	"encoding/hex"
)

// hash returns the hex encoded hash of b
func hash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
