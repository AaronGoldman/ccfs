package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
)

type commit struct {
	listHash   []byte
	versionstr string
	hkid       []byte
	signature  []byte //131 byte max
}

func NewCommit(listHash []byte, hkid []byte) (c commit) {
	c.listHash = listHash
	c.versionstr = GenerateVersion()
	c.hkid = hkid
	c.signature = commitSign(c.listHash, c.versionstr, c.hkid)
	return
}

func commitSign(listHash []byte, versionstr string, hkid []byte) (signature []byte) {
	ObjectHash := commitRefHash(listHash, versionstr, hkid)
	prikey, err := getPrivateKeyForHkid(hkid)
	r, s, err := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	if err != nil {
		panic(err)
	}

	signature = elliptic.Marshal(prikey.PublicKey.Curve, r, s)
	return
}

func (c commit) Verifiy() (ret bool) {
	fmt.Print(1)
	ObjectHash := commitRefHash(c.listHash, c.versionstr, c.hkid)
	fmt.Print(2)
	pubkey := getPiblicKeyForHkid(c.hkid)
	fmt.Print(pubkey)
	r, s := elliptic.Unmarshal(pubkey.Curve, c.signature)
	fmt.Print(4)
	ret = ecdsa.Verify(pubkey, ObjectHash, r, s)
	fmt.Print(5)
	return
}

func commitRefHash(listHash []byte, versionstr string, hkid []byte) (ObjectHash []byte) {
	var h hash.Hash = sha256.New()
	h.Write(listHash[:])
	h.Write([]byte(versionstr))
	h.Write(hkid[:])
	ObjectHash = h.Sum(make([]byte, 0))
	return
}

func GenerateCommit(list []byte, key *ecdsa.PrivateKey) (Commit string) {
	listHashBytes := GenerateObjectHash(list)
	listHashStr := hex.EncodeToString(listHashBytes)
	versionstr := GenerateVersion()
	signature := GenerateSignature(key, listHashBytes)
	hkidstr := GenerateHKID(key)
	return fmt.Sprintf("%s,%s,%s,%s", listHashStr,
		versionstr, hkidstr, signature)
}

// func GenerateCommit(list []byte, key *ecdsa.PrivateKey) (Commit string) {
//	listHashBytes := GenerateObjectHash(list)
//	listHashStr := hex.EncodeToString(listHashBytes)
//	versionstr := GenerateVersion()
//	signature := GenerateSignature(key, listHashBytes)
//	hkidstr := GenerateHKID(key)
//	return fmt.Sprintf("%s,%s,%s,%s", listHashStr,
//		versionstr, hkidstr, signature)
//}
