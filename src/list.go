package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

type entry struct {
	Hash        []byte
	TypeString  string
	nameSegment string
}

func (e entry) String() string {
	return fmt.Sprintf("%s,%s,%s", hex.EncodeToString(e.Hash),
		e.TypeString, e.nameSegment)
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
		s = s + fmt.Sprintf("%s\n", element.String())
	}
	return s[:len(s)-1]
}

func (l list) hash_for_namesegment(namesegment string) (string, []byte) {
	for _, element := range l {
		if strings.EqualFold(element.nameSegment, namesegment) {
			return element.TypeString, element.Hash
		}
	}
	return "null", nil
}

func NewList(hash []byte, typestring string, nameSegment string) list {
	e := entry{hash, typestring, nameSegment}
	return list{e}
}

func NewListFromBytes(listbytes []byte) (newlist list) {
	listEntries := strings.Split(string(listbytes), "\n")
	entries := []entry{}
	cols := []string{}
	for _, element := range listEntries {
		cols = strings.Split(element, ",")
		entryHash, _ := hex.DecodeString(cols[0])
		entryTypeString := cols[1]
		entryNameSegment := cols[2]
		entries = append(entries, entry{entryHash, entryTypeString,
			entryNameSegment})
	}
	newlist = list(entries)
	return
}

func GetList(hash []byte) (l list, err error) {
	listbytes, err := GetBlob(hash)
	l = NewListFromBytes(listbytes)
	return
}
