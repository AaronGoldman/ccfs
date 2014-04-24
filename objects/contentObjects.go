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
		entryHash, _ := hex.DecodeString(cols[0])
		entryTypeString := cols[1]
		var entryHID HID
		if entryTypeString == "blob" || entryTypeString == "list" {
			entryHID = HCID(entryHash)
		} else if entryTypeString == "commit" || entryTypeString == "tag" {
			entryHID = HKID(entryHash)
		} else {
			log.Fatal()
		}
		entryNameSegment := cols[2]
		l[entryNameSegment] = entry{entryHID, entryTypeString}
	}
	return l, err
}

type Commit struct {
	listHash  HCID
	version   int64
	parent    HCID
	hkid      HKID
	signature []byte //131 byte max
}

func (c Commit) ListHash() HCID {
	return c.listHash
}

func (c Commit) Version() int64 {
	return c.version
}

func (c Commit) Parent() HCID {
	return c.parent
}

func (c Commit) Hkid() HKID {
	return c.hkid
}

func (c Commit) Signature() []byte {
	return c.signature
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
		hex.EncodeToString(c.listHash),
		c.version,
		hex.EncodeToString(c.parent),
		hex.EncodeToString(c.hkid),
		hex.EncodeToString(c.signature))
}

func (c Commit) Verify() bool {
	ObjectHash := c.genCommitHash(c.listHash, c.version, c.parent, c.hkid)
	pubkey := ecdsa.PublicKey(geterPoster.getPiblicKeyForHkid(c.hkid))
	if pubkey.Curve == nil || pubkey.X == nil || pubkey.Y == nil {
		return false
	}
	r, s := elliptic.Unmarshal(pubkey.Curve, c.signature)
	//log.Println(pubkey, " pubkey\n", ObjectHash, " ObjectHash\n", r, " r\n", s, "s")
	if r.BitLen() == 0 || s.BitLen() == 0 {
		return false
	}
	return ecdsa.Verify(&pubkey, ObjectHash, r, s)
}

func (c Commit) Update(listHash HCID) Commit {
	c.parent = c.Hash()
	c.version = time.Now().UnixNano()
	//c.hkid = c.hkid
	c.listHash = listHash
	c.signature = c.commitSign(c.listHash, c.version, c.parent, c.hkid)
	return c
}

func (c Commit) commitSign(listHash []byte, version int64, parent HCID, hkid []byte) (signature []byte) {
	ObjectHash := c.genCommitHash(listHash, version, parent, hkid)
	prikey, err := geterPoster.getPrivateKeyForHkid(hkid)
	ecdsaprikey := ecdsa.PrivateKey(*prikey)
	r, s, err := ecdsa.Sign(rand.Reader, &ecdsaprikey, ObjectHash)
	if err != nil {
		log.Panic(err)
	}
	signature = elliptic.Marshal(prikey.PublicKey.Curve, r, s)
	return
}

func (c Commit) genCommitHash(listHash []byte, version int64, parent HCID, hkid []byte) (
	ObjectHash []byte) {
	var h hash.Hash = sha256.New()
	h.Write([]byte(fmt.Sprintf("%s,\n%d,\n%s,\n%s",
		hex.EncodeToString(listHash),
		version,
		parent,
		hex.EncodeToString(hkid))))
	ObjectHash = h.Sum(nil)
	return
}

func NewCommit(listHash []byte, hkid HKID) (c Commit) {
	c.listHash = listHash
	c.version = time.Now().UnixNano()
	c.hkid = hkid
	c.parent = sha256.New().Sum(nil)
	c.signature = c.commitSign(c.listHash, c.version, c.parent, c.hkid)
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
	listHash, _ := hex.DecodeString(commitStrings[0])
	version, _ := strconv.ParseInt(commitStrings[1], 10, 64)
	parent, _ := hex.DecodeString(commitStrings[2])
	cHkid, _ := hex.DecodeString(commitStrings[3])
	signature, _ := hex.DecodeString(commitStrings[4])
	//var h hash.Hash = sha256.New()
	c = Commit{listHash, version, parent, cHkid, signature}
	return
}

type Tag struct {
	HashBytes   HID
	TypeString  string
	NameSegment string
	Version     int64
	hkid        HKID
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
	return fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s,\n%s",
		hex.EncodeToString(t.HashBytes.Bytes()),
		t.TypeString,
		t.NameSegment,
		t.Version,
		hex.EncodeToString(t.hkid),
		hex.EncodeToString(t.Signature))
}

func (t Tag) Hkid() HKID {
	return t.hkid
}

func (t Tag) Verify() bool {
	if t.Hkid == nil {
		return false
	}
	tPublicKey := ecdsa.PublicKey(geterPoster.getPiblicKeyForHkid(t.hkid))
	r, s := elliptic.Unmarshal(elliptic.P521(), t.Signature)
	ObjectHash := t.genTagHash(t.HashBytes, t.TypeString, t.NameSegment,
		t.Version, t.hkid)
	if r.BitLen() == 0 || s.BitLen() == 0 {
		return false
	}
	return ecdsa.Verify(&tPublicKey, ObjectHash, r, s)
}

func (t Tag) Update(hashBytes HID, typeString string) Tag {
	t.HashBytes = hashBytes
	t.TypeString = typeString
	//t.nameSegment = t.nameSegment
	t.Version = time.Now().UnixNano()
	//t.hkid = t.hkid
	prikey, err := geterPoster.getPrivateKeyForHkid(t.hkid)
	if err != nil {
		log.Panic("You don't seem to own this Domain")
	}
	ObjectHash := t.genTagHash(
		t.HashBytes,
		t.TypeString,
		t.NameSegment,
		t.Version,
		t.hkid)
	ecdsaprikey := ecdsa.PrivateKey(*prikey)
	r, s, _ := ecdsa.Sign(rand.Reader, &ecdsaprikey, ObjectHash)
	t.Signature = elliptic.Marshal(elliptic.P521(), r, s)
	return t
}

func NewTag(HashBytes HID, TypeString string,
	nameSegment string, hkid HKID) Tag {
	prikey, _ := geterPoster.getPrivateKeyForHkid(hkid)
	version := time.Now().UnixNano()
	ObjectHash := Tag{}.genTagHash(HashBytes, TypeString, nameSegment, version, hkid)
	ecdsaprikey := ecdsa.PrivateKey(*prikey)
	r, s, _ := ecdsa.Sign(rand.Reader, &ecdsaprikey, ObjectHash)
	signature := elliptic.Marshal(elliptic.P521(), r, s)
	t := Tag{HashBytes,
		TypeString,
		nameSegment,
		version,
		hkid,
		signature}
	return t
}

func (t Tag) genTagHash(Hash HID, TypeString string, nameSegment string,
	version int64, hkid []byte) []byte {
	var h hash.Hash = sha256.New()
	h.Write([]byte(fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s",
		hex.EncodeToString(Hash.Bytes()),
		TypeString,
		nameSegment,
		version,
		hex.EncodeToString(hkid))))
	return h.Sum(nil)
}

func TagFromBytes(bytes []byte) (t Tag, err error) {
	//build object
	tagStrings := strings.Split(string(bytes), ",\n")
	if len(tagStrings) != 6 {
		return t, fmt.Errorf("Could not parse tag bytes")
	}
	tagHashBytes, _ := hex.DecodeString(tagStrings[0])
	tagTypeString := tagStrings[1]
	var tagHID HID
	if tagTypeString == "blob" || tagTypeString == "list" {
		tagHID = HCID(tagHashBytes)
	} else if tagTypeString == "commit" || tagTypeString == "tag" {
		tagHID = HKID(tagHashBytes)
	} else {
		log.Fatal()
	}
	tagNameSegment := tagStrings[2]
	tagVersion, _ := strconv.ParseInt(tagStrings[3], 10, 64)
	tagHkid, _ := hex.DecodeString(tagStrings[4])
	tagSignature, _ := hex.DecodeString(tagStrings[5])
	t = Tag{tagHID, tagTypeString, tagNameSegment, tagVersion, tagHkid,
		tagSignature}
	return
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