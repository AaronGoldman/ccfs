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

type HCID []byte

func (hcid HCID) Bytes() []byte {
	return hcid
}

func (hcid HCID) Hex() string {
	return hex.EncodeToString(hcid)
}

func (hcid HCID) String() string {
	return hcid.Hex()
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

type HKID []byte

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

func HkidFromHex(s string) (HKID, error) {
	dabytes, err := hex.DecodeString(s)
	if err == nil {
		return HKID(dabytes), err
	}
	return nil, err
}

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

type GeterPoster struct {
	getPiblicKeyForHkid  func(hkid HKID) PublicKey
	getPrivateKeyForHkid func(hkid HKID) (k *PrivateKey, err error)
	PostKey              func(p *PrivateKey) error
	PostBlob             func(b Blob) error
}

var geterPoster GeterPoster

func RegisterGeterPoster(
	getPiblicKeyForHkid func(hkid HKID) PublicKey,
	getPrivateKeyForHkid func(hkid HKID) (k *PrivateKey, err error),
	PostKey func(p *PrivateKey) error,
	PostBlob func(b Blob) error,
) {
	geterPoster.getPiblicKeyForHkid = getPiblicKeyForHkid
	geterPoster.getPrivateKeyForHkid = getPrivateKeyForHkid
	geterPoster.PostKey = PostKey
	geterPoster.PostBlob = PostBlob
}
