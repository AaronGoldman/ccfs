package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func commitgentest() {
	priv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	fmt.Printf("%v", GenerateCommit([]byte("testing"), priv))
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
