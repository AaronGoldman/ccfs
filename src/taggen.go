package main

//import (
//	"crypto/ecdsa"
//	"encoding/hex"
//	"fmt"
//)

///*TAG
//ObjectHash(HEX),
//ObjectType,
//NameSegment(url escaped),
//Version,
//HKID,
//Signature(r s)
//*/

////string:=QueryEscape(s string)
////string, error:=QueryUnescape(s string)

//func GenerateTag(blob []byte, objectType string, nameSegment string,
//	key *ecdsa.PrivateKey) (tag string) {
//	objectHashBytes := GenerateObjectHash(blob)
//	objectHashStr := hex.EncodeToString(objectHashBytes)
//	objectType = GenerateObjectType(objectType)
//	nameSegment = GenerateNameSegment(nameSegment)
//	versionstr := GenerateVersion()
//	signature := GenerateSignature(key, objectHashBytes)
//	hkidstr := GenerateHKID(key)

//	tag = fmt.Sprintf("%s,\n%s,\n%s,\n%s,\n%s,\n%s", objectHashStr,
//		objectType, nameSegment, versionstr, hkidstr, signature)
//	return tag
//}
