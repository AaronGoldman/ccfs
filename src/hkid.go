package main

import (
	"encoding/hex"
)

type HKID HCID

func (hkid HKID) Hex() string {
	return hex.EncodeToString(hkid)
}
