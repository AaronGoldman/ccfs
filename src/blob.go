package main

import (
	"crypto/sha256"
	"hash"
)

type blob []byte

func (b blob) Hash() []byte {
	var h hash.Hash = sha256.New()
	h.Write(b)
	return h.Sum(make([]byte, 0))
}

//func ([]byte b) Hash() []byte {
//	var h hash.Hash = sha256.New()
//	h.Write(b)
//	return h.Sum(make([]byte, 0))
//}

//func NewBlob(inbyte []byte)(blob){
//	outblob := new(blob)
//	outblob = blob(inbyte)
//	return outblob
//}
