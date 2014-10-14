//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package objects

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

//HID is a Hash IDentifier
type HID interface {
	Byteser
	Hexer
}

//Byteser is a object that can be exported as a slice of byte
type Byteser interface {
	Bytes() []byte
}

//Hexer is a object that can be exported as a string of hex
type Hexer interface {
	Hex() string
}

//HCID is the type to represent content by it's cryptographic hash
type HCID []byte

//Bytes returns a slise of byte representing the hash
func (hcid HCID) Bytes() []byte {
	return hcid
}

//Hex returns a string containing a hexadecimal representation if the hash
func (hcid HCID) Hex() string {
	return hex.EncodeToString(hcid)
}

//String returns a string containing a hexadecimal representation if the hash
func (hcid HCID) String() string {
	return hcid.Hex()
}

//HcidFromHex return an HCID if the string contains 64 digits of Hex
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

//HKID is the hash of a public key stored in a slice of byte
type HKID []byte

//Bytes reterns the HKID in the form of a slice of byte
func (hkid HKID) Bytes() []byte {
	return hkid
}

//Hex reterns the HKID in the form of a hexidesimal string.
func (hkid HKID) Hex() string {
	return hex.EncodeToString(hkid)
}

func (hkid HKID) String() string {
	return hkid.Hex()
}

//GenHKID returns a new HKID and posts the key to the blob store
func GenHKID() HKID {
	privkey := KeyGen()
	err := geterPoster.PostKey(privkey)
	if err != nil {
		log.Fatalf("Failed to persist Privet Key")
	}
	err = geterPoster.PostBlob(elliptic.Marshal(privkey.PublicKey.Curve,
		privkey.PublicKey.X, privkey.PublicKey.Y))
	if err != nil {
		log.Fatalf("Failed to post Public Key")
	}
	return privkey.Hkid()
}

//HkidFromHex gives you an HKID form a hex string
func HkidFromHex(s string) (HKID, error) {
	dabytes, err := hex.DecodeString(s)
	if err == nil {
		return HKID(dabytes), err
	}
	return nil, err
}

//HkidFromD builds a HKID using a number in a big int
func HkidFromD(D big.Int) HKID {
	priv, err := PrivteKeyFromD(D)
	key := elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y)
	hkid := priv.Hkid()
	err = geterPoster.PostKey(priv) //store privet key
	if err != nil {
		log.Panic(err)
	}
	err = geterPoster.PostBlob(key) //store public key
	if err != nil {
		log.Panic(err)
	}
	return hkid
}

//HkidFromDString builds a HKID using a number in base in a string
func HkidFromDString(str string, base int) HKID {
	D, success := new(big.Int).SetString(str, base)
	if !success {
		log.Panic(nil)
	}
	return HkidFromD(*D)
}

type parents []HCID

func (p parents) String() string {
	parentString := ""
	for _, pHCID := range p {
		parentString = parentString + "," + pHCID.Hex()
	}
	return parentString[1:]
}

//GeterPoster is a struct of functions for use with objects
type GeterPoster struct {
	getPublicKeyForHkid  func(hkid HKID) PublicKey
	getPrivateKeyForHkid func(hkid HKID) (k *PrivateKey, err error)
	PostKey              func(p *PrivateKey) error
	PostBlob             func(b Blob) error
}

var geterPoster GeterPoster

//RegisterGeterPoster ads a given GeterPoster for use with objects
func RegisterGeterPoster(
	getPublicKeyForHkid func(hkid HKID) PublicKey,
	getPrivateKeyForHkid func(hkid HKID) (k *PrivateKey, err error),
	PostKey func(p *PrivateKey) error,
	PostBlob func(b Blob) error,
) {
	geterPoster.getPublicKeyForHkid = getPublicKeyForHkid
	geterPoster.getPrivateKeyForHkid = getPrivateKeyForHkid
	geterPoster.PostKey = PostKey
	geterPoster.PostBlob = PostBlob
}
