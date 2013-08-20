package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	//"log"
	"sort"
	"strings"
)

type entry struct {
	Hash       HID
	TypeString string
}

type list map[string]entry

func (l list) add(nameSegment string, hash HID, typeString string) list {
	l[nameSegment] = entry{hash, typeString}
	return l
}

func (l list) hash_for_namesegment(namesegment string) (string, Hexer) {
	objectHash := l[namesegment].Hash
	typeString := l[namesegment].TypeString
	return typeString, objectHash
}

func (l list) String() string {
	var keys []string
	for key := range l {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	s := ""
	for _, k := range keys {
		s = s + fmt.Sprintf("%s,%s,%s\n", l[k].Hash.Hex(), l[k].TypeString, k)
	}
	return s[:len(s)-1]
}

func (l list) Bytes() []byte {
	return []byte(l.String())
}

func (l list) Hash() []byte {
	var h hash.Hash = sha256.New()
	h.Write(l.Bytes())
	return h.Sum(nil)
}

func NewList(objectHash []byte, typestring string, nameSegment string) list {
	l := make(list)
	l[nameSegment] = entry{objectHash, typestring}
	return l
}

func NewListFromBytes(listbytes []byte) (newlist list) {
	l := make(list)
	listEntries := strings.Split(string(listbytes), "\n")
	cols := []string{}
	for _, element := range listEntries {
		cols = strings.Split(element, ",")
		//log.Print(cols)
		entryHash, _ := hex.DecodeString(cols[0])
		entryTypeString := cols[1]
		entryNameSegment := cols[2]
		l[entryNameSegment] = entry{entryHash, entryTypeString}
	}
	return l
}

func GetList(objectHash HCID) (l list, err error) {
	listbytes, err := GetBlob(objectHash)
	if len(listbytes) == 0 {
		return nil, err
	}
	l = NewListFromBytes(listbytes)
	return
}

func PostList(l list) (err error) {
	return PostBlob(blob(l.Bytes()))
}
