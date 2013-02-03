package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
)

func GenerateCommit(list []byte, key *ecdsa.PrivateKey) (Commit string) {
	listHashBytes := GenerateObjectHash(list)
	listHashStr := hex.EncodeToString(listHashBytes)
	versionstr := GenerateVersion()
	signature := GenerateSignature(key, listHashBytes)
	hkidstr := GenerateHKID(key)
	return fmt.Sprintf("%s,%s,%s,%s", listHashStr,
		versionstr, hkidstr, signature)
}
