package main

import (
	"encoding/hex"
)

type HCID []byte

func (hcid HCID) Hex() string {
	return hex.EncodeToString(hcid)
}
