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
	"encoding/gob"
	"math/big"
"os"
)

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

type privkey struct{
	pubkey
	D *big.Int
}

type pubkey struct{
	X, Y *big.Int
}

func SaveKey(priv *ecdsa.PrivateKey){
	filepath := fmt.Sprintf("../keys/%s",GenerateHKID(priv))
	fmt.Print(*priv)
	fo, err := os.Create(filepath)

	if err != nil {
		fmt.Printf("%v", err)
		return
		//panic(err)
	}
	enc := gob.NewEncoder(fo)
	err = enc.Encode(privkey{pubkey{priv.PublicKey.X, priv.PublicKey.Y}, priv.D})
	if err != nil {
        fmt.Printf("\n%v", err)
    }
    fo.Close()
}

func LoadKey(hkid string)(priv *ecdsa.PrivateKey){
filepath := fmt.Sprintf("../keys/%s",hkid)
fmt.Print(filepath)
fi, err := os.Create(filepath)

	if err != nil {
		fmt.Printf("\n%v", err)
		return
		//panic(err)
	}
	defer fi.Close()
dec := gob.NewDecoder(fi)
var prikey privkey
err = dec.Decode(&prikey)
if err != nil {
        panic(fmt.Sprintf("\n%v", err))
    }
priv = &ecdsa.PrivateKey{ecdsa.PublicKey{elliptic.P521(), prikey.pubkey.X, prikey.pubkey.Y}, prikey.D}
fmt.Printf("%v\n",priv)
return priv
}
