//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package objects

import (
	"crypto/sha256"
	"fmt"
	"log"
	"sort"
	"strings"
)

type entry struct {
	Hash       HID
	TypeString string
}

type List map[string]entry

func (l List) Add(nameSegment string, hash HID, typeString string) List {
	l[nameSegment] = entry{hash, typeString}
	return l
}

func (l List) Remove(nameSegment string) List {
	delete(l, nameSegment)
	return l
}

func (l List) HashForNamesegment(namesegment string) (string, HID) {
	objectHash := l[namesegment].Hash
	typeString := l[namesegment].TypeString
	return typeString, objectHash
}

func (l List) String() string {
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

func (l List) Log() {
	log.Printf(
		"List %s\n-----BEGIN LIST-------\n%q\n-------END LIST-------",
		l.Hash(),
		l,
	)
}

func (l List) Bytes() []byte {
	return []byte(l.String())
}

func (l List) Hash() HCID {
	h := sha256.New()
	h.Write(l.Bytes())
	return h.Sum(nil)
}

func NewList(objectHash HID, typestring string, nameSegment string) List {
	l := make(List)
	l[nameSegment] = entry{objectHash, typestring}
	return l
}

func ListFromBytes(listbytes []byte) (newlist List, err error) {
	l := make(List)
	listEntries := strings.Split(string(listbytes), "\n")
	cols := []string{}
	for _, element := range listEntries {
		cols = strings.Split(element, ",")
		if len(cols) != 3 {
			return newlist, fmt.Errorf("Could not parse list")
		}
		entryTypeString := cols[1]
		var entryHID HID
		switch entryTypeString {
		case "blob", "list":
			entryHID, err = HcidFromHex(cols[0])
			if err != nil {
				return nil, err
			}
		case "commit", "tag":
			entryHID, err = HkidFromHex(cols[0])
			if err != nil {
				return nil, err
			}
		default:
			log.Fatalf("Unrecognised type: %s", entryTypeString)
		}

		entryNameSegment := cols[2]
		l[entryNameSegment] = entry{entryHID, entryTypeString}
	}
	return l, err
}
