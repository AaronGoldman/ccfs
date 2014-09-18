package objects

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//Tag is a individualy vertioned reference to a content object
type Tag struct {
	HashBytes   HID
	TypeString  string
	NameSegment string
	Version     int64
	Parents     parents
	Hkid        HKID
	Signature   []byte
}

//Hash caculates the HCID of the tag
func (t Tag) Hash() HCID {
	h := sha256.New()
	h.Write(t.Bytes())
	return h.Sum(nil)
}

//Bytes builds text representation of the tag as a slice of byte
func (t Tag) Bytes() []byte {
	return []byte(t.String())
}

//Sting builds text representation of the tag
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

//Log sends the tag to the log with a header and footer
func (t Tag) Log() {
	log.Printf(
		"Posting tag %s\n-----BEGIN TAG--------\n%q\n-------END TAG--------",
		t.Hash(),
		t,
	)
}

//Verify that the signature is valid
func (t Tag) Verify() bool {
	if t.Hkid == nil {
		return false
	}
	tPublicKey := ecdsa.PublicKey(geterPoster.getPublicKeyForHkid(t.Hkid))
	r, s := elliptic.Unmarshal(elliptic.P521(), t.Signature)
	ObjectHash := t.genTagHash(t.HashBytes, t.TypeString, t.NameSegment,
		t.Version, t.Parents, t.Hkid)
	if r.BitLen() == 0 || s.BitLen() == 0 {
		return false
	}
	return ecdsa.Verify(&tPublicKey, ObjectHash, r, s)
}

//Update the named content
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

//Delete the named content
func (t Tag) Delete() Tag {
	t = NewTag(
		Blob{}.Hash(),     //HashBytes HID
		"nab",             //TypeString string
		t.NameSegment,     //nameSegment string
		parents{t.Hash()}, //tparent parents
		t.Hkid,            //hkid HKID
	)
	return t
}

//Merge updates the content and add extra the tags to the parents
func (t Tag) Merge(tags []Tag, hashBytes HID, typeString string) Tag {
	tParents := parents{t.Hash()}
	for _, pTag := range tags {
		tParents = append(tParents, pTag.Hash())
	}
	t = NewTag(
		hashBytes,     //HashBytes HID
		typeString,    //TypeString string
		t.NameSegment, //nameSegment string
		tParents,      //tparent parents
		t.Hkid,        //hkid HKID
	)
	return t
}

//NewTag build a new tag with the inital content
func NewTag(
	HashBytes HID,
	TypeString string,
	nameSegment string,
	tparent parents,
	hkid HKID,
) Tag {
	prikey, _ := geterPoster.getPrivateKeyForHkid(hkid)
	version := time.Now().UnixNano()
	if tparent == nil {
		tparent = parents{Blob{}.Hash()}
	}
	ObjectHash := Tag{}.genTagHash(
		HashBytes,
		TypeString,
		nameSegment,
		version,
		tparent,
		hkid,
	)
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
	h := sha256.New()
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

//TagFromBytes parses a slice of byte to a tag
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
	switch tagTypeString {
	case "blob", "list":
		tagHID, err = HcidFromHex(tagStrings[0])
		if err != nil {
			return
		}
	case "commit", "tag":
		tagHID, err = HkidFromHex(tagStrings[0])
		if err != nil {
			return
		}
	default:
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
