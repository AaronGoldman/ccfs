package objects

import (
	"encoding/hex"
	"fmt"
)

type HCID []byte

func (hcid HCID) Hex() string {
	return hex.EncodeToString(hcid)
}

func (hcid HCID) Bytes() []byte {
	return hcid
}

type HKID []byte

//Hex reterns the HKID in the form of a hexidesimal string.
func (hkid HKID) Hex() string {
	return hex.EncodeToString(hkid)
}

func (hkid HKID) Bytes() []byte {
	return hkid
}

func (hkid HKID) String() string {
	return hkid.Hex()
}

func HkidFromHex(s string) (HKID, error) {
	dabytes, err := hex.DecodeString(s)
	if err == nil {
		return HKID(dabytes), err
	}
	return nil, err
}

func HcidFromHex(s string) (HCID, error) {
	dabytes, err := hex.DecodeString(s)
	if err == nil {
		return HCID(dabytes), err
	}
	if len(s) != 64 {
		return nil, fmt.Errorf("HEX not 64 digits")
	}
	return nil, err
}

func (hcid HCID) String() string {
	return hcid.Hex()
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
