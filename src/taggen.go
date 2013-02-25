package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
)

/*TAG
ObjectHash(HEX),
ObjectType,
NameSegment(url escaped),
Version,
HKID,
Signature(r s)
*/

//string:=QueryEscape(s string)
//string, error:=QueryUnescape(s string)

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
	objectHashStr := hex.EncodeToString(t.HashBytes)
	nameSegment := GenerateNameSegment(t.nameSegment)
	signature := string(t.signature)
	hkidstr := string(t.hkid)
	tagstring := fmt.Sprintf("%s,\n%s,\n%s,\n%s,\n%s,\n%s", objectHashStr,
		t.TypeString, nameSegment, t.versionstr, hkidstr, signature)
	return tagstring
}

func GenerateTag(blob []byte, objectType string, nameSegment string,
	key *ecdsa.PrivateKey) (tag string) {
	objectHashBytes := GenerateObjectHash(blob)
	objectHashStr := hex.EncodeToString(objectHashBytes)
	objectType = GenerateObjectType(objectType)
	nameSegment = GenerateNameSegment(nameSegment)
	versionstr := GenerateVersion()
	signature := GenerateSignature(key, objectHashBytes)
	hkidstr := GenerateHKID(key)

	tag = fmt.Sprintf("%s,\n%s,\n%s,\n%s,\n%s,\n%s", objectHashStr,
		objectType, nameSegment, versionstr, hkidstr, signature)
	return tag
}
