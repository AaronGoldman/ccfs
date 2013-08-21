package main

import (
	"encoding/hex"
)

type HCID []byte

func (hcid HCID) Hex() string {
	return hex.EncodeToString(hcid)
	//HID(hcid).Hex()
}

func (hcid HCID) Bytes() []byte {
	return hcid
}

type HKID []byte

func (hkid HKID) Hex() string {
	return hex.EncodeToString(hkid)
	//HID(hkid).Hex()
}

func (hkid HKID) Bytes() []byte {
	return hkid
}

type HID interface {
	Byteser
	Hexer
}
type Byteser interface {
	Bytes() []byte
}

type Hexer interface {
	Hex() string
}
