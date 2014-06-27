//Criptograficly curated content objects
//
//This packages contains the contnet objects
//    and functions for working with them
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
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
)

type GeterPoster struct {
	getPiblicKeyForHkid  func(hkid HKID) PublicKey
	getPrivateKeyForHkid func(hkid HKID) (k *PrivateKey, err error)
	PostKey              func(p *PrivateKey) error
	PostBlob             func(b Blob) error
}

var geterPoster GeterPoster

func RegisterGeterPoster(
	getPiblicKeyForHkid func(hkid HKID) PublicKey,
	getPrivateKeyForHkid func(hkid HKID) (k *PrivateKey, err error),
	PostKey func(p *PrivateKey) error,
	PostBlob func(b Blob) error,
) {
	geterPoster.getPiblicKeyForHkid = getPiblicKeyForHkid
	geterPoster.getPrivateKeyForHkid = getPrivateKeyForHkid
	geterPoster.PostKey = PostKey
	geterPoster.PostBlob = PostBlob
}

type Blob []byte

func (b Blob) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(b)
	return HCID(h.Sum(make([]byte, 0)))
}

func (b Blob) Bytes() []byte {
	return []byte(b)
}

type entry struct {
	Hash       HID
	TypeString string
}

type List map[string]entry

func (l List) Add(nameSegment string, hash HID, typeString string) List {
	l[nameSegment] = entry{hash, typeString}
	return l
}

func (l List) Remove(nameSegment string) List{
	delete(l, nameSegment)
	return l
}

func (l List) Hash_for_namesegment(namesegment string) (string, HID) {
	objectHash := l[namesegment].Hash
	typeString := l[namesegment].TypeString
	return typeString, objectHash
}

func (l List) String() string {
	var keys []string
	for key := range l {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	s := ""
	for _, k := range keys {
		s = s + fmt.Sprintf("%s,%s,%s\n", l[k].Hash.Hex(), l[k].TypeString, k)
	}
	return s[:len(s)-1]
}

func (l List) Bytes() []byte {
	return []byte(l.String())
}

func (l List) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(l.Bytes())
	return h.Sum(nil)
}

func NewList(objectHash HID, typestring string, nameSegment string) List {
	l := make(List)
	l[nameSegment] = entry{objectHash, typestring}
	return l
}

func ListFromBytes(listbytes []byte) (newlist List, err error) {
	l := make(List)
	listEntries := strings.Split(string(listbytes), "\n")
	cols := []string{}
	for _, element := range listEntries {
		cols = strings.Split(element, ",")
		//log.Print(cols)
		//entryHash, err := hex.DecodeString(cols[0])
		//if err != nil {
		//	return nil, err
		//}
		if len(cols) != 3 {
			return newlist, fmt.Errorf("Could not parse list")
		}
		entryTypeString := cols[1]
		var entryHID HID
		if entryTypeString == "blob" || entryTypeString == "list" {
			entryHID, err = HcidFromHex(cols[0])
			if err != nil {
				return nil, err
			}
		} else if entryTypeString == "commit" || entryTypeString == "tag" {
			entryHID, err = HkidFromHex(cols[0])
			if err != nil {
				return nil, err
			}
		} else {
			log.Fatal()
		}
		entryNameSegment := cols[2]
		l[entryNameSegment] = entry{entryHID, entryTypeString}
	}
	return l, err
}

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

//func InitCommit() HKID {
//	privkey := KeyGen()
//	hkid := privkey.Hkid()
//	PostCommit(NewCommit(sha256.New().Sum(nil), hkid))
//	return hkid
//}

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

type Tag struct {
	HashBytes   HID
	TypeString  string
	NameSegment string
	Version     int64
	Parents     parents
	Hkid        HKID
	Signature   []byte
}

func (t Tag) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(t.Bytes())
	return h.Sum(nil)
}

func (t Tag) Bytes() []byte {
	return []byte(t.String())
}

func (t Tag) String() string {

	return fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s,\n%s,\n%s",
		t.HashBytes.Hex(),
		t.TypeString,
		t.NameSegment,
		t.Version,
		t.Parents,
		t.Hkid.Hex(),
		hex.EncodeToString(t.Signature),
	)
}

func (t Tag) Verify() bool {
	if t.Hkid == nil {
		return false
	}
	tPublicKey := ecdsa.PublicKey(geterPoster.getPiblicKeyForHkid(t.Hkid))
	r, s := elliptic.Unmarshal(elliptic.P521(), t.Signature)
	ObjectHash := t.genTagHash(t.HashBytes, t.TypeString, t.NameSegment,
		t.Version, t.Parents, t.Hkid)
	if r.BitLen() == 0 || s.BitLen() == 0 {
		return false
	}
	return ecdsa.Verify(&tPublicKey, ObjectHash, r, s)
}

func (t Tag) Update(hashBytes HID, typeString string) Tag {
	t.Parents = parents{t.Hash()}
	t.HashBytes = hashBytes
	t.TypeString = typeString
	//t.nameSegment = t.nameSegment
	t.Version = time.Now().UnixNano()
	//t.hkid = t.hkid
	prikey, err := geterPoster.getPrivateKeyForHkid(t.Hkid)
	if err != nil {
		log.Panic("You don't seem to own this Domain")
	}

	ObjectHash := t.genTagHash(
		t.HashBytes,
		t.TypeString,
		t.NameSegment,
		t.Version,
		t.Parents,
		t.Hkid,
	)
	ecdsaprikey := ecdsa.PrivateKey(*prikey)
	r, s, _ := ecdsa.Sign(rand.Reader, &ecdsaprikey, ObjectHash)
	t.Signature = elliptic.Marshal(elliptic.P521(), r, s)
	return t
}

func (t Tag) Merge(tags []Tag, hashBytes HID, typeString string) Tag {
	t.Parents = parents{t.Hash()}
	for _, pTag := range tags {
		t.Parents = append(t.Parents, pTag.Hash())
	}
	t.HashBytes = hashBytes
	t.TypeString = typeString
	//t.nameSegment = t.nameSegment
	t.Version = time.Now().UnixNano()
	//t.hkid = t.hkid
	prikey, err := geterPoster.getPrivateKeyForHkid(t.Hkid)
	if err != nil {
		log.Panic("You don't seem to own this Domain")
	}

	ObjectHash := t.genTagHash(
		t.HashBytes,
		t.TypeString,
		t.NameSegment,
		t.Version,
		t.Parents,
		t.Hkid,
	)
	ecdsaprikey := ecdsa.PrivateKey(*prikey)
	r, s, _ := ecdsa.Sign(rand.Reader, &ecdsaprikey, ObjectHash)
	t.Signature = elliptic.Marshal(elliptic.P521(), r, s)
	return t
}

func NewTag(HashBytes HID, TypeString string,
	nameSegment string, tparent parents, hkid HKID) Tag {
	prikey, _ := geterPoster.getPrivateKeyForHkid(hkid)
	version := time.Now().UnixNano()
	if tparent == nil {
		tparent = parents{Blob{}.Hash()}
	}
	ObjectHash := Tag{}.genTagHash(HashBytes, TypeString, nameSegment, version,
		tparent, hkid)

	ecdsaprikey := ecdsa.PrivateKey(*prikey)
	r, s, _ := ecdsa.Sign(rand.Reader, &ecdsaprikey, ObjectHash)
	signature := elliptic.Marshal(elliptic.P521(), r, s)
	t := Tag{HashBytes,
		TypeString,
		nameSegment,
		version,
		tparent,
		hkid,
		signature}
	return t
}

func (t Tag) genTagHash(taghash HID, TypeString string, nameSegment string,
	version int64, tparent parents, hkid HKID) []byte {
	var h hash.Hash = sha256.New()
	h.Write(
		[]byte(fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s,\n%s",
			taghash.Hex(),
			TypeString,
			nameSegment,
			version,
			tparent,
			hkid,
		)),
	)
	return h.Sum(nil)
}

func TagFromBytes(bytes []byte) (t Tag, err error) {
	//build object
	tagStrings := strings.Split(string(bytes), ",\n")
	if len(tagStrings) != 7 {
		return t, fmt.Errorf("Could not parse tag bytes")
	}
	//tagHashBytes, err := hex.DecodeString(tagStrings[0])
	//if err != nil {
	//	return nil, err
	//}
	tagTypeString := tagStrings[1]
	var tagHID HID
	if tagTypeString == "blob" || tagTypeString == "list" {
		tagHID, err = HcidFromHex(tagStrings[0])
		if err != nil {
			return
		}
	} else if tagTypeString == "commit" || tagTypeString == "tag" {
		tagHID, err = HkidFromHex(tagStrings[0])
		if err != nil {
			return
		}
	} else {
		return t, fmt.Errorf("Could not parse tag")
	}
	tagNameSegment := tagStrings[2]
	tagVersion, err := strconv.ParseInt(tagStrings[3], 10, 64)
	if err != nil {
		return
	}

	parentSplit := strings.Split(tagStrings[4], ",")
	parsedParents := parents{}
	for _, singlParentString := range parentSplit {
		parsedHCID, err1 := HcidFromHex(singlParentString)
		if err1 != nil {
			return t, err1
		}
		parsedParents = append(parsedParents, parsedHCID)
	}

	tagHkid, err := hex.DecodeString(tagStrings[5])
	if err != nil {
		return
	}
	tagSignature, err := hex.DecodeString(tagStrings[6])
	if err != nil {
		return
	}
	t = Tag{
		tagHID,
		tagTypeString,
		tagNameSegment,
		tagVersion,
		parsedParents,
		tagHkid,
		tagSignature,
	}
	return
}

type parents []HCID

func (p parents) String() string {
	parentString := ""
	for _, pHCID := range p {
		parentString = parentString + "," + pHCID.Hex()
	}
	return parentString[1:]
}

func GenHKID() HKID {
	privkey := KeyGen()
	err := geterPoster.PostKey(privkey)
	if err != nil {
		log.Fatalf("Failed to persist Privet Key")
	}
	err = geterPoster.PostBlob(elliptic.Marshal(privkey.PublicKey.Curve,
		privkey.PublicKey.X, privkey.PublicKey.Y))
	if err != nil {
		log.Fatalf("Failed to post Public Key")
	}
	return privkey.Hkid()
}

func HkidFromDString(str string, base int) HKID {
	D, success := new(big.Int).SetString(str, base)
	if !success {
		log.Panic(nil)
	}
	return HkidFromD(*D)
}

func HkidFromD(D big.Int) HKID {
	priv, err := PrivteKeyFromD(D)
	key := elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y)
	hkid := priv.Hkid()
	err = geterPoster.PostKey(priv) //store privet key
	if err != nil {
		log.Panic(err)
	}
	err = geterPoster.PostBlob(key) //store public key
	if err != nil {
		log.Panic(err)
	}
	return hkid
}

//func GenerateObjectType(objectType string) (objectTypestr string) {
//	return objectType
//}
//func GenerateNameSegment(nameSegment string) (nameSegmentstr string) {
//	return url.QueryEscape(nameSegment)
//}
//func GenerateVersion() (versionstr string) {
//	return strconv.FormatInt(time.Now().UnixNano(), 10)
//}
//func GenerateSignature(prikey *ecdsa.PrivateKey, ObjectHash []byte) (signature string) {
//	r, s, err := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
//	if err != nil {
//		log.Panic(err)
//	}
//	return fmt.Sprintf("%v %v", r, s)
//}
