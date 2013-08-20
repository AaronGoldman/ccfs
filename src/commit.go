package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"log"
	"strconv"
	"strings"
	"time"
)

type commit struct {
	listHash  HCID
	version   int64
	parent    HCID
	hkid      HKID
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
	return fmt.Sprintf("%s,\n%d,\n%s,\n%s,\n%s",
		hex.EncodeToString(c.listHash),
		c.version,
		hex.EncodeToString(c.parent),
		hex.EncodeToString(c.hkid),
		hex.EncodeToString(c.signature))
}

func (c commit) Hkid() []byte {
	return c.hkid
}

func (c commit) Verifiy() bool {
	ObjectHash := genCommitHash(c.listHash, c.version, c.hkid)
	pubkey := ecdsa.PublicKey(getPiblicKeyForHkid(c.hkid))
	r, s := elliptic.Unmarshal(pubkey.Curve, c.signature)
	//log.Println(pubkey, " pubkey\n", ObjectHash, " ObjectHash\n", r, " r\n", s, "s")
	return ecdsa.Verify(&pubkey, ObjectHash, r, s)
}

func (c commit) Update() commit {
	return c
}

func commitSign(listHash []byte, version int64, hkid []byte) (signature []byte) {
	ObjectHash := genCommitHash(listHash, version, hkid)
	prikey, err := getPrivateKeyForHkid(hkid)
	r, s, err := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	if err != nil {
		log.Panic(err)
	}
	signature = elliptic.Marshal(prikey.PublicKey.Curve, r, s)
	return
}

func genCommitHash(listHash []byte, version int64, hkid []byte) (
	ObjectHash []byte) {
	var h hash.Hash = sha256.New()
	h.Write([]byte(fmt.Sprintf("%s,\n%d,\n%s",
		hex.EncodeToString(listHash),
		version,
		//parent,
		hex.EncodeToString(hkid))))
	ObjectHash = h.Sum(nil)
	return
}

func NewCommit(listHash []byte, hkid HKID) (c commit) {
	c.listHash = listHash
	c.version = time.Now().UnixNano()
	c.hkid = hkid
	c.parent = sha256.New().Sum(nil)
	c.signature = commitSign(c.listHash, c.version, c.hkid)
	return
}

func InitCommit() HKID {
	privkey := KeyGen()
	hkid := GenerateHKID(privkey)
	PostCommit(NewCommit(sha256.New().Sum(nil), hkid))
	return hkid
}

func CommitFromBytes(bytes []byte) (c commit, err error) {
	//build object
	commitStrings := strings.Split(string(bytes), ",\n")
	listHash, _ := hex.DecodeString(commitStrings[0])
	version, _ := strconv.ParseInt(commitStrings[1], 10, 64)
	parent, _ := hex.DecodeString(commitStrings[2])
	cHkid, _ := hex.DecodeString(commitStrings[3])
	signature, _ := hex.DecodeString(commitStrings[4])
	//var h hash.Hash = sha256.New()
	c = commit{listHash, version, parent, cHkid, signature}
	return
}
