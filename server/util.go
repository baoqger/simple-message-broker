package server

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func genId() string {
	u := make([]byte, 16)
	io.ReadFull(rand.Reader, u)
	return hex.EncodeToString(u)
}
