package main

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"fmt"
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
	return h.Sum(make([]byte, 0))
}

func (l list) Bytes() []byte {
	return []byte(l.String())
}

func (l list) String() string {
	for index, element := range l {
		element.String()
	}
}

func (e entry) String() string {
	return fmt.Sprintf("%s,%s,%s", hex.EncodeToString(e.Hash), e.TypeString, e.nameSegment)
}

func NewList() list {
	return list{[]entry{
			[{[]byte("test"),
			"blob",
			"place holder"}]}
	}
}

func GenerateList(blobs [][]byte, objectTypes []string,
	nameSegment []string) (list string) {
	return list
}
