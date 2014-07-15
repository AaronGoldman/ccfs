package objects

import (
	"crypto/sha256"
	"hash"
)

type Blob []byte

func (b Blob) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(b)
	return HCID(h.Sum(make([]byte, 0)))
}

func (b Blob) Bytes() []byte {
	return []byte(b)
}
