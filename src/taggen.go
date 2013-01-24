package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
)

/*TAG
ObjectHash(HEX),
ObjectType,
NameSegment(url escaped),
Version,
Signature(r s),
HKID
*/

//string:=QueryEscape(s string)
//string, error:=QueryUnescape(s string)

func main() {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	fmt.Printf("priv:%v\n err:%v\n", priv, err)
	hashed := []byte("testing")

	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)
	fmt.Printf("sig\nr:%v\ns:%v\n err:%v\n", r, s, err)

	valid := ecdsa.Verify(&priv.PublicKey, hashed, r, s)
	invalid := ecdsa.Verify(&priv.PublicKey, []byte("fail testing"), r, s)

	fmt.Printf("valid:%v in:%v\n", valid, invalid)
}

func GenerateTag(blob []byte, objectType string, nameSegment string,
	version int, key *ecdsa.PrivateKey) (tag string) {
	objectHashBytes := GenerateObjectHash(blob)
	objectHashStr := hex.EncodeToString(objectHashBytes)
	objectType = GenerateObjectType(objectType)
	nameSegment = GenerateNameSegment(nameSegment)
	versionstr := GenerateVersion(version)
	signature := GenerateSignature(key, objectHashBytes)
	hkidstr := GenerateHKID(key)

	tag = fmt.Sprintf("%s,%s,%s,%s,%s,%s", objectHashStr, objectType, nameSegment,
		versionstr, signature, hkidstr)
	return tag
}

func GenerateObjectHash(blob []byte) (objectHash []byte) {
	return []byte("ToDo")
}
func GenerateObjectType(objectType string) (objectTypestr string) {
	return objectType
}
func GenerateNameSegment(nameSegment string) (nameSegmentstr string) {
	return url.QueryEscape(nameSegment)
}
func GenerateVersion(version int) (versionstr string) {
	return strconv.Itoa(version)
}
func GenerateSignature(prikey *ecdsa.PrivateKey, ObjectHash []byte) (signature string) {
	return "ToDo"
}
func GenerateHKID(prikey *ecdsa.PrivateKey) (hkid string) {
	return "ToDo"
}
