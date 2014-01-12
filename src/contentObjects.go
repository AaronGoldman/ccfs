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
	//"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type blob []byte

func (b blob) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(b)
	return HCID(h.Sum(make([]byte, 0)))
}

func (b blob) Bytes() []byte {
	return []byte(b)
}

type entry struct {
	Hash       HID
	TypeString string
}

type list map[string]entry

func (l list) add(nameSegment string, hash HID, typeString string) list {
	l[nameSegment] = entry{hash, typeString}
	return l
}

func (l list) hash_for_namesegment(namesegment string) (string, HID) {
	objectHash := l[namesegment].Hash
	typeString := l[namesegment].TypeString
	return typeString, objectHash
}

func (l list) String() string {
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

func (l list) Bytes() []byte {
	return []byte(l.String())
}

func (l list) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(l.Bytes())
	return h.Sum(nil)
}

func NewList(objectHash HID, typestring string, nameSegment string) list {
	l := make(list)
	l[nameSegment] = entry{objectHash, typestring}
	return l
}

func NewListFromBytes(listbytes []byte) (newlist list) {
	l := make(list)
	listEntries := strings.Split(string(listbytes), "\n")
	cols := []string{}
	for _, element := range listEntries {
		cols = strings.Split(element, ",")
		//log.Print(cols)
		entryHash, _ := hex.DecodeString(cols[0])
		entryTypeString := cols[1]
		entryNameSegment := cols[2]
		l[entryNameSegment] = entry{HCID(entryHash), entryTypeString}
	}
	return l
}

func GetList(objectHash HCID) (l list, err error) {
	listbytes, err := GetBlob(objectHash)
	if len(listbytes) == 0 {
		return nil, err
	}
	l = NewListFromBytes(listbytes)
	return
}

func PostList(l list) (err error) {
	return PostBlob(blob(l.Bytes()))
}

type commit struct {
	listHash  HCID
	version   int64
	parent    HCID
	hkid      HKID
	signature []byte //131 byte max
}

func (c commit) Hash() HCID {
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
	ObjectHash := genCommitHash(c.listHash, c.version, c.parent, c.hkid)
	pubkey := ecdsa.PublicKey(getPiblicKeyForHkid(c.hkid))
	r, s := elliptic.Unmarshal(pubkey.Curve, c.signature)
	//log.Println(pubkey, " pubkey\n", ObjectHash, " ObjectHash\n", r, " r\n", s, "s")
	return ecdsa.Verify(&pubkey, ObjectHash, r, s)
}

func (c commit) Update(listHash HCID) commit {
	c.parent = c.Hash()
	c.version = time.Now().UnixNano()
	//c.hkid = c.hkid
	c.listHash = listHash
	c.signature = commitSign(c.listHash, c.version, c.parent, c.hkid)
	return c
}

func commitSign(listHash []byte, version int64, parent HCID, hkid []byte) (signature []byte) {
	ObjectHash := genCommitHash(listHash, version, parent, hkid)
	prikey, err := getPrivateKeyForHkid(hkid)
	r, s, err := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	if err != nil {
		log.Panic(err)
	}
	signature = elliptic.Marshal(prikey.PublicKey.Curve, r, s)
	return
}

func genCommitHash(listHash []byte, version int64, parent HCID, hkid []byte) (
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

func NewCommit(listHash []byte, hkid HKID) (c commit) {
	c.listHash = listHash
	c.version = time.Now().UnixNano()
	c.hkid = hkid
	c.parent = sha256.New().Sum(nil)
	c.signature = commitSign(c.listHash, c.version, c.parent, c.hkid)
	return
}

func InitCommit() HKID {
	privkey := KeyGen()
	hkid := privkey.Hkid()
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

type tag struct {
	HashBytes   HCID
	TypeString  string
	nameSegment string
	version     int64
	hkid        HKID
	signature   []byte
}

func (t tag) Hash() HCID {
	var h hash.Hash = sha256.New()
	h.Write(t.Bytes())
	return h.Sum(nil)
}

func (t tag) Bytes() []byte {
	return []byte(t.String())
}

func (t tag) String() string {
	return fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s,\n%s",
		hex.EncodeToString(t.HashBytes),
		t.TypeString,
		t.nameSegment,
		t.version,
		hex.EncodeToString(t.hkid),
		hex.EncodeToString(t.signature))
}

func (t tag) Hkid() HKID {
	return t.hkid
}

func (t tag) Verifiy() bool {
	tPublicKey := ecdsa.PublicKey(getPiblicKeyForHkid(t.hkid))
	r, s := elliptic.Unmarshal(elliptic.P521(), t.signature)
	ObjectHash := genTagHash(t.HashBytes, t.TypeString, t.nameSegment,
		t.version, t.hkid)
	return ecdsa.Verify(&tPublicKey, ObjectHash, r, s)
}

func (t tag) Update(hashBytes HCID) tag {
	t.HashBytes = hashBytes
	//t.TypeString = typeString
	//t.nameSegment = t.nameSegment
	t.version = time.Now().UnixNano()
	//t.hkid = t.hkid
	prikey, err := getPrivateKeyForHkid(t.hkid)
	if err != nil {
		log.Panic("You don't seem to own this Domain")
	}
	ObjectHash := genTagHash(
		t.HashBytes.Bytes(),
		t.TypeString,
		t.nameSegment,
		t.version,
		t.hkid)
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	t.signature = elliptic.Marshal(elliptic.P521(), r, s)
	return t
}

func NewTag(HashBytes HID, TypeString string,
	nameSegment string, hkid HKID) tag {
	prikey, _ := getPrivateKeyForHkid(hkid)
	version := time.Now().UnixNano()
	ObjectHash := genTagHash(HashBytes.Bytes(), TypeString, nameSegment, version, hkid)
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, ObjectHash)
	signature := elliptic.Marshal(elliptic.P521(), r, s)
	t := tag{HashBytes.Bytes(),
		TypeString,
		nameSegment,
		version,
		hkid,
		signature}
	return t
}

func genTagHash(HashBytes []byte, TypeString string, nameSegment string,
	version int64, hkid []byte) []byte {
	var h hash.Hash = sha256.New()
	h.Write([]byte(fmt.Sprintf("%s,\n%s,\n%s,\n%d,\n%s",
		hex.EncodeToString(HashBytes),
		TypeString,
		nameSegment,
		version,
		hex.EncodeToString(hkid))))
	return h.Sum(nil)
}

func TagFromBytes(bytes []byte) (t tag, err error) {
	//build object
	tagStrings := strings.Split(string(bytes), ",\n")
	tagHashBytes, _ := hex.DecodeString(tagStrings[0])
	tagTypeString := tagStrings[1]
	tagNameSegment := tagStrings[2]
	tagVersion, _ := strconv.ParseInt(tagStrings[3], 10, 64)
	tagHkid, _ := hex.DecodeString(tagStrings[4])
	tagSignature, _ := hex.DecodeString(tagStrings[5])
	t = tag{tagHashBytes, tagTypeString, tagNameSegment, tagVersion, tagHkid,
		tagSignature}
	return
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
