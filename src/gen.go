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

func SaveKey(priv *ecdsa.PrivateKey){
	filepath := fmt.Sprintf("../keys/%s",GenerateHKID(priv))
	fo, err := os.Create(filepath)

	if err != nil {
		fmt.Printf("%v", err)
		return
		//panic(err)
	}
	defer fo.Close()
	enc := gob.NewEncoder(fo)
	enc.Encode(priv)
	
}

func LoadKey(hkid string)(priv *ecdsa.PrivateKey){
filepath := fmt.Sprintf("../keys/%s",hkid)
fi, err := os.Create(filepath)

	if err != nil {
		fmt.Printf("%v", err)
		return
		//panic(err)
	}
	defer fi.Close()
dec := gob.NewDecoder(fi)
//var priv ecdsa.PrivateKey
dec.Decode(priv)
return priv
}
