package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"time"
)

type commit struct {
	listHash  []byte
	version   int64
	hkid      []byte
	signature []byte //131 byte max
}

func (c commit) Hash() []byte {
	var h hash.Hash = sha256.New()
	h.Write(c.Bytes())
	return h.Sum(nil)
}

func (c commit) Bytes() []byte {
	return []byte(c.String())
}

func (c commit) String() string {
	return fmt.Sprintf("%s,\n%d,\n%s,\n%s",
		hex.EncodeToString(c.listHash),
		c.version,
		hex.EncodeToString(c.hkid),
		hex.EncodeToString(c.signature))
}

func (c commit) Hkid() []byte {
	return c.hkid
}

func genCommitHash(listHash []byte, version int64, hkid []byte) (ObjectHash []byte) {
	var h hash.Hash = sha256.New()
	h.Write([]byte(fmt.Sprintf("%s,\n%d,\n%s",
		hex.EncodeToString(listHash),
		version,
		hex.EncodeToString(hkid))))
	ObjectHash = h.Sum(nil)
	return
}

func commitSign(listHash []byte, version int64, hkid []byte) (signature []byte) {
	ObjectHash := genCommitHash(listHash, version, hkid)
	prikey, err := getPrivateKeyForHkid(hkid)
	r, s, err := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	if err != nil {
		panic(err)
	}
	signature = elliptic.Marshal(prikey.PublicKey.Curve, r, s)
	return
}

func (c commit) Verifiy() bool {
	ObjectHash := genCommitHash(c.listHash, c.version, c.hkid)
	pubkey := getPiblicKeyForHkid(c.hkid)
	r, s := elliptic.Unmarshal(pubkey.Curve, c.signature)
	return ecdsa.Verify(pubkey, ObjectHash, r, s)
}
func NewCommit(listHash []byte, hkid []byte) (c commit) {
	c.listHash = listHash
	c.version = time.Now().UnixNano()
	c.hkid = hkid
	c.signature = commitSign(c.listHash, c.version, c.hkid)
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
