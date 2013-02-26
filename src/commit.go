package main

import (
//"crypto/ecdsa"
//"encoding/hex"
//"fmt"
)

type commit struct {
	listHash   [32]byte
	versionstr string
	hkid       [32]byte
	signature  []byte //131 byte max
}

func (t commit) Verifiy() bool {
	return false
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
