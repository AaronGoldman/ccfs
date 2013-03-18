package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
)

type Tag struct {
	HashBytes   []byte
	TypeString  string
	nameSegment string
	versionstr  int64
	signature   []byte //(r, s *big.Int)
	hkid        []byte
	//func Marshal(curve Curve, x, y *big.Int) []byte
	//func Unmarshal(curve Curve, data []byte) (x, y *big.Int)
	//elliptic.Marshal(prikey.PublicKey.Curve,prikey.PublicKey.X,prikey.PublicKey.Y)
}

func (t Tag) String() string {
	objectHashStr := hex.EncodeToString(t.HashBytes[:])
	nameSegment := GenerateNameSegment(t.nameSegment)
	signature := string(t.signature)
	hkidstr := string(t.hkid[:])
	tagstring := fmt.Sprintf("%s,\n%s,\n%s,\n%s,\n%s,\n%s", objectHashStr,
		t.TypeString, nameSegment, t.versionstr, hkidstr, signature)
	return tagstring
}

func (t Tag) Verifiy() bool {
	PublicKey := getPiblicKeyForHkid(t.hkid)
	r, s := elliptic.Unmarshal(elliptic.P521(), t.signature)
	hashed := []byte("testing") //place holder
	return ecdsa.Verify(PublicKey, hashed, r, s)
}

func NewTag(HashBytes   []byte, TypeString  string, 
	nameSegment string, hkid []byte){
	
}
