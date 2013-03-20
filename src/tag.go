package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"time"
)

type Tag struct {
	HashBytes   []byte
	TypeString  string
	nameSegment string
	version     int64
	hkid        []byte
	signature   []byte //Marshal(r, s *big.Int)
	//func Marshal(curve Curve, x, y *big.Int) []byte
	//func Unmarshal(curve Curve, data []byte) (x, y *big.Int)
	//elliptic.Marshal(prikey.PublicKey.Curve,prikey.PublicKey.X,prikey.PublicKey.Y)
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
	return fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s",
		hex.EncodeToString(t.HashBytes),
		t.TypeString,
		t.nameSegment,
		t.version,
		hex.EncodeToString(t.hkid),
		hex.EncodeToString(t.signature))
}

func (t Tag) Verifiy() bool {
	PublicKey := getPiblicKeyForHkid(t.hkid)
	r, s := elliptic.Unmarshal(elliptic.P521(), t.signature)
	hashed := []byte("testing") //place holder
	return ecdsa.Verify(PublicKey, hashed, r, s)
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

func NewTag(HashBytes []byte, TypeString string,
	nameSegment string, hkid []byte) Tag {
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
