package objects

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

type Commit struct {
	ListHash  HCID
	Version   int64
	Parents   parents
	Hkid      HKID
	Signature []byte //131 byte max
}

func (c Commit) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(c.Bytes())
	return h.Sum(nil)
}

func (c Commit) Bytes() []byte {
	return []byte(c.String())
}

func (c Commit) String() string {
	return fmt.Sprintf("%s,\n%d,\n%s,\n%s,\n%s",
		c.ListHash.Hex(),
		c.Version,
		c.Parents,
		c.Hkid.Hex(),
		hex.EncodeToString(c.Signature))
}

func (c Commit) Verify() bool {
	ObjectHash := c.genCommitHash(c.ListHash, c.Version, c.Parents, c.Hkid)
	pubkey := ecdsa.PublicKey(geterPoster.getPiblicKeyForHkid(c.Hkid))
	if pubkey.Curve == nil || pubkey.X == nil || pubkey.Y == nil {
		return false
	}
	r, s := elliptic.Unmarshal(pubkey.Curve, c.Signature)
	//log.Println(pubkey, " pubkey\n", ObjectHash, " ObjectHash\n", r, " r\n", s, "s")
	if r.BitLen() == 0 || s.BitLen() == 0 {
		return false
	}
	return ecdsa.Verify(&pubkey, ObjectHash, r, s)
}

func (c Commit) Update(listHash HCID) Commit {
	c.Parents = parents{c.Hash()}
	c.Version = time.Now().UnixNano()
	//c.Hkid = c.Hkid
	c.ListHash = listHash
	c.Signature = c.commitSign(c.ListHash, c.Version, c.Parents, c.Hkid)
	return c
}

func (c Commit) Merge(pCommits []Commit, listHash HCID) Commit {
	c.Parents = parents{c.Hash()}
	for _, pCommit := range pCommits {
		c.Parents = append(c.Parents, pCommit.Hash())
	}
	c.Version = time.Now().UnixNano()
	//c.Hkid = c.Hkid
	c.ListHash = listHash
	c.Signature = c.commitSign(c.ListHash, c.Version, c.Parents, c.Hkid)
	return c
}

func (c Commit) commitSign(listHash []byte, version int64, cparents parents, hkid []byte) (signature []byte) {
	ObjectHash := c.genCommitHash(listHash, version, cparents, hkid)
	prikey, err := geterPoster.getPrivateKeyForHkid(hkid)
	ecdsaprikey := ecdsa.PrivateKey(*prikey)
	r, s, err := ecdsa.Sign(rand.Reader, &ecdsaprikey, ObjectHash)
	if err != nil {
		log.Panic(err)
	}
	signature = elliptic.Marshal(prikey.PublicKey.Curve, r, s)
	return
}

func (c Commit) genCommitHash(
	listHash HCID,
	version int64,
	cparents parents,
	hkid HKID,
) (ObjectHash []byte) {
	var h hash.Hash = sha256.New()
	h.Write([]byte(fmt.Sprintf("%s,\n%d,\n%s,\n%s",
		listHash,
		version,
		cparents,
		hkid,
	)))
	ObjectHash = h.Sum(nil)
	return
}

func NewCommit(listHash HCID, hkid HKID) (c Commit) {
	c.ListHash = listHash
	c.Version = time.Now().UnixNano()
	c.Hkid = hkid
	c.Parents = []HCID{sha256.New().Sum(nil)}
	c.Signature = c.commitSign(c.ListHash, c.Version, c.Parents, c.Hkid)
	return
}

func CommitFromBytes(bytes []byte) (c Commit, err error) {
	//build object
	commitStrings := strings.Split(string(bytes), ",\n")
	if len(commitStrings) != 5 {
		return c, fmt.Errorf("Could not parse commit bytes")
	}
	listHash, err := hex.DecodeString(commitStrings[0])
	if err != nil {
		return
	}
	version, err := strconv.ParseInt(commitStrings[1], 10, 64)
	if err != nil {
		return
	}
	parentSplit := strings.Split(commitStrings[2], ",")
	parsedParents := parents{}
	for _, singlParentString := range parentSplit {
		parsedHCID, err1 := HcidFromHex(singlParentString)
		if err1 != nil {
			return c, err1
		}
		parsedParents = append(parsedParents, parsedHCID)
	}

	cHkid, err := hex.DecodeString(commitStrings[3])
	if err != nil {
		return
	}
	signature, err := hex.DecodeString(commitStrings[4])
	if err != nil {
		return
	}
	c = Commit{listHash, version, parsedParents, cHkid, signature}
	return
}
