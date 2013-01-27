package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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
