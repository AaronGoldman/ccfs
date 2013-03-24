package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
)

type entry struct {
	Hash        []byte
	TypeString  string
	nameSegment string
}

type list []entry

func (l list) Hash() []byte {
	var h hash.Hash = sha256.New()
	h.Write(l.Bytes())
	return h.Sum(nil)
}

func (l list) Bytes() []byte {
	return []byte(l.String())
}

func (l list) String() string {
	s := ""
	for _, element := range l {
		s = fmt.Sprintf("%s\n", element.String())
	}
	return s[:len(s)-1]
}

func (e entry) String() string {
	return fmt.Sprintf("%s,%s,%s", hex.EncodeToString(e.Hash),
		e.TypeString, e.nameSegment)
}

func NewList(hash []byte, typestring string, nameSegment string) list {
	e := entry{hash, typestring, nameSegment}
	return list{e}
}

//func GenerateList(blobs [][]byte, objectTypes []string,
//	nameSegment []string) (list string) {
//	return list
//}
