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
	SaveKey(priv)
	hkid := GenerateHKID(priv)
	reconstructedkey := LoadKey("11c18133c89f71ec4edbe31a137a2b789935f5e87b4b8a91acad493937dda73a")
	reconstructedhkid := GenerateHKID(reconstructedkey)
	fmt.Printf("\n%v\n%v", hkid, reconstructedhkid)
	//fmt.Printf("%v", GenerateCommit([]byte("testing"), priv))
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
