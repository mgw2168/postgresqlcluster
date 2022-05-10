package pkg

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1Str(s string) string {
	r := sha1.Sum([]byte(s))
	foo := hex.EncodeToString(r[:])

	return string([]byte(foo)[:10])
}
