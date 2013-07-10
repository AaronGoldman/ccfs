package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	HashBytes   []byte
	TypeString  string
	nameSegment string
	version     int64
	hkid        []byte
	signature   []byte
}

func (t Tag) Hash() []byte {
	var h hash.Hash = sha256.New()
	h.Write(t.Bytes())
	return h.Sum(nil)
}

func (t Tag) Bytes() []byte {
	return []byte(t.String())
}

func (t Tag) String() string {
	return fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s,\n%s",
		hex.EncodeToString(t.HashBytes),
		t.TypeString,
		t.nameSegment,
		t.version,
		hex.EncodeToString(t.hkid),
		hex.EncodeToString(t.signature))
}

func (t Tag) Hkid() []byte {
	return t.hkid
}

func (t Tag) Verifiy() bool {
	PublicKey := getPiblicKeyForHkid(t.hkid)
	r, s := elliptic.Unmarshal(elliptic.P521(), t.signature)
	ObjectHash := genTagHash(t.HashBytes, t.TypeString, t.nameSegment,
		t.version, t.hkid)
	return ecdsa.Verify(PublicKey, ObjectHash, r, s)
}

func (t Tag) Update() Tag {
	return t
}

func NewTag(HashBytes HCID, TypeString string,
	nameSegment string, hkid HKID) Tag {
	prikey, _ := getPrivateKeyForHkid(hkid)
	version := time.Now().UnixNano()
	ObjectHash := genTagHash(HashBytes, TypeString, nameSegment, version, hkid)
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	signature := elliptic.Marshal(elliptic.P521(), r, s)
	t := Tag{HashBytes,
		TypeString,
		nameSegment,
		version,
		hkid,
		signature}
	return t
}

func genTagHash(HashBytes []byte, TypeString string, nameSegment string,
	version int64, hkid []byte) []byte {
	var h hash.Hash = sha256.New()
	h.Write([]byte(fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s",
		hex.EncodeToString(HashBytes),
		TypeString,
		nameSegment,
		version,
		hex.EncodeToString(hkid))))
	return h.Sum(nil)
}

func TagFromBytes(bytes []byte) (t Tag, err error) {
	//build object
	tagStrings := strings.Split(string(bytes), ",\n")
	tagHashBytes, _ := hex.DecodeString(tagStrings[0])
	tagTypeString := tagStrings[1]
	tagNameSegment := tagStrings[2]
	tagVersion, _ := strconv.ParseInt(tagStrings[3], 10, 64)
	tagHkid, _ := hex.DecodeString(tagStrings[4])
	tagSignature, _ := hex.DecodeString(tagStrings[5])
	t = Tag{tagHashBytes, tagTypeString, tagNameSegment, tagVersion, tagHkid,
		tagSignature}
	return
}
