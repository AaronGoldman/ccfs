package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"log"
	"net/url"
	"strconv"
	"time"
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
	r, s, err := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf("%v %v", r, s)
}

func GenerateHKID(prikey *ecdsa.PrivateKey) (hkid HKID) {
	var h hash.Hash = sha256.New()
	h.Write(elliptic.Marshal(
		prikey.PublicKey.Curve,
		prikey.PublicKey.X,
		prikey.PublicKey.Y))
	return h.Sum(make([]byte, 0))
}
