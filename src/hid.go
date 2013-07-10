package main

import (
	"encoding/hex"
)

type HID []byte

func (hid HID) Hex() string {
	return hex.EncodeToString(hid)
}

func (hid HID) Bytes() []byte {
	return []byte(hid)
}

type HCID []byte

func (hcid HCID) Hex() string {
	return HID(hcid).Hex()
	//hex.EncodeToString(hcid)
}

func (hcid HCID) Bytes() []byte {
	return HID(hcid).Bytes()
}

type HKID HCID

func (hkid HKID) Hex() string {
	return HID(hkid).Hex()
	//hex.EncodeToString(hkid)
}

func (hkid HKID) Bytes() []byte {
	return HID(hkid).Bytes()
}

type Byteser interface {
	Bytes() []byte
}

type Hexer interface {
	Hex() string
}
