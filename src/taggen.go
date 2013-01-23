package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net/url"
)

/*TAG
ObjectHash(HEX),
ObjectType,
NameSegment(url escaped),
Version,
Signature,
HKID
*/

//string:=QueryEscape(s string)
//string, error:=QueryUnescape(s string)

func main() {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	fmt.Printf("priv:%v\n err:%v\n", priv, err)
	hashed := []byte("testing")

	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)
	fmt.Printf("r:%v\n s:%v\n err:%v\n", r, s, err)

	valid := ecdsa.Verify(&priv.PublicKey, hashed, r, s)
	invalid := ecdsa.Verify(&priv.PublicKey, []byte("fail testing"), r, s)

	fmt.Printf("valid:%v in:%v\n", valid, invalid)
}
