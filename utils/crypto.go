package utils

import (
	"crypto/sha256"
	"fmt"
)

func StringToSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x\n", bs)
}
