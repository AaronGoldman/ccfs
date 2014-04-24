package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"log"
	"strings"
)

//Get retrieves the content objects using HID and path
func Get(objecthash objects.HID, path string) (b objects.Blob, err error) {
	typeString := "commit"
	err = nil
	nameSegments := []string{"", path}
	for {
		if len(nameSegments) > 1 {
			nameSegments = strings.SplitN(nameSegments[1], "/", 2)
		} else {
			nameSegments = []string{""}
		}
		//log.Printf("\n\tPath: %s\n\tType: %v\n\tobjecthash: %v\n",
		//	nameSegments, typeString, objecthash)
		switch typeString {
		case "blob":
			b, err = GetBlob(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			return b, err
		case "list":
			var l objects.List
			l, err = GetList(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.Hash_for_namesegment(nameSegments[0])
			if objecthash == nil && nameSegments[0] != "" {
				err = fmt.Errorf("Blob not found")
			}
			b = l.Bytes()
		case "tag":
			var t objects.Tag
			//if nameSegments[0] == "" {
			//	log.Printf("\n\tNo Path\n")
			//}
			t, err = GetTag(objecthash.(objects.HKID), nameSegments[0])
			if err != nil {
				//log.Printf("\n\t%v\n", err)
				return nil, err
			}
			if !t.Verify() {
				return nil, fmt.Errorf("Tag Verifiy Failed")
			}
			typeString = t.TypeString
			objecthash = t.HashBytes
			b = t.Bytes()
		case "commit":
			var c objects.Commit
			c, err = GetCommit(objecthash.Bytes())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			if !c.Verify() {
				return nil, fmt.Errorf("Commit Verifiy Failed")
			}
			var l objects.List
			l, err = GetList(c.ListHash())
			if err != nil {
				log.Printf("\n\t%v\n", err)
			}
			typeString, objecthash = l.Hash_for_namesegment(nameSegments[0])
			if objecthash == nil && nameSegments[0] != "" {
				err = fmt.Errorf("Blob not found")
			}
			//if err != nil {
			//	log.Printf("%v\n", err)
			//}
			b = l.Bytes()
		default:
			log.Printf("\n\t%v\n", err)
			panic(err)
		}
		//if len(nameSegments) == 1 && typeString != "blob" {
		if objecthash == nil {
			return b, err
		}
	}
}

func GetList(objectHash objects.HCID) (l objects.List, err error) {
	listbytes, err := GetBlob(objectHash)
	if len(listbytes) == 0 {
		return nil, err
	}
	l, err = objects.ListFromBytes(listbytes)
	return
}

//getPiblicKeyForHkid uses the lookup services to get a public key for an hkid
func GetPiblicKeyForHkid(hkid objects.HKID) objects.PublicKey {
	marshaledKey, err := GetBlob(objects.HCID(hkid))
	if err != nil {
		return objects.PublicKey{}
	}
	curve := elliptic.P521()
	x, y := elliptic.Unmarshal(elliptic.P521(), marshaledKey)
	pubKey := ecdsa.PublicKey{
		Curve: curve, //elliptic.Curve
		X:     x,     //X *big.Int
		Y:     y}     //Y *big.Int
	return objects.PublicKey(pubKey)
}

//getPrivateKeyForHkid uses the lookup services to get a private key for an hkid
func GetPrivateKeyForHkid(hkid objects.HKID) (k *objects.PrivateKey, err error) {
	k, err = GetKey(hkid)
	return k, err
}