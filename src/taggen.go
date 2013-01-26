package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"net/url"
	"strconv"
	"time"
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

func taggentest() {
	priv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//test(priv)
	fmt.Printf("%v\n", GenerateTag([]byte("testing"), "blob", "test", priv))
}

func test(priv *ecdsa.PrivateKey) {
	mar := elliptic.Marshal(
		priv.PublicKey.Curve,
		priv.PublicKey.X,
		priv.PublicKey.Y)
	x, y := elliptic.Unmarshal(elliptic.P521(), mar)
	maredPublicKey := new(ecdsa.PublicKey)
	maredPublicKey.Curve = elliptic.P521()
	maredPublicKey.X = x
	maredPublicKey.Y = y

	hashed := []byte("testing")
	r, s, _ := ecdsa.Sign(rand.Reader, priv, hashed)
	valid := ecdsa.Verify(maredPublicKey, hashed, r, s)
	invalid := ecdsa.Verify(maredPublicKey, []byte("fail testing"), r, s)
	fmt.Printf("valid:%v in:%v marsize:%v bits\n\n\n", valid, invalid, len(mar)*8)
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

func GenerateObjectHash(blob []byte) (objectHash []byte) {
	var h hash.Hash = sha256.New()
	h.Write(blob)
	return h.Sum(make([]byte, 0))
}
func GenerateObjectType(objectType string) (objectTypestr string) {
	return objectType
}
func GenerateNameSegment(nameSegment string) (nameSegmentstr string) {
	return url.QueryEscape(nameSegment)
}
func GenerateVersion() (versionstr string) {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
func GenerateSignature(prikey *ecdsa.PrivateKey, ObjectHash []byte) (signature string) {
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	return fmt.Sprintf("%v %v", r, s)
}
func GenerateHKID(prikey *ecdsa.PrivateKey) (hkid string) {
	var h hash.Hash = sha256.New()
	h.Write(elliptic.Marshal(
		prikey.PublicKey.Curve,
		prikey.PublicKey.X,
		prikey.PublicKey.Y))
	return fmt.Sprintf("%v", hex.EncodeToString(h.Sum(make([]byte, 0))))
}
