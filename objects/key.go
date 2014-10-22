//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package objects

//D *big.Int is all the key matereal save as bin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

//PrivateKey wrapper around the ecdsa.PrivateKey type
type PrivateKey ecdsa.PrivateKey

//PublicKey wrapper around the ecdsa.PrivateKey type
type PublicKey ecdsa.PublicKey

//Bytes returns the marshaled public key as a slice of byte.
func (p PublicKey) Bytes() []byte {
	return elliptic.Marshal(p.Curve, p.X, p.Y)
}

//Hkid returns the hkid for the public key.
func (p PublicKey) Hkid() HKID {
	h := sha256.New()
	h.Write(p.Bytes())
	return h.Sum(make([]byte, 0))
}

//Hkid returns the hkid that for private key.
//this is the hcid of your public key
func (p PrivateKey) Hkid() HKID {
	h := sha256.New()
	h.Write(elliptic.Marshal(p.PublicKey.Curve, p.PublicKey.X, p.PublicKey.Y))
	return h.Sum(make([]byte, 0))
}

//Hash returns the hcid for the PrivateKey
func (p PrivateKey) Hash() HCID {
	h := sha256.New()
	h.Write(p.Bytes())
	return h.Sum(nil)
}

//Verify returns true if PrivateKey and PrivateKey are a pair.
func (p PrivateKey) Verify() bool {
	ObjectHash := make([]byte, 32)
	_, err := rand.Read(ObjectHash)
	if err != nil {
		fmt.Println("error:", err)
	}
	prikey := ecdsa.PrivateKey(p)
	r, s, err := ecdsa.Sign(rand.Reader, &prikey, ObjectHash)
	return ecdsa.Verify(&p.PublicKey, ObjectHash, r, s)
}

//Bytes returns the marshaled public key as a slice of byte.
func (p PrivateKey) Bytes() []byte {
	return ecdsa.PrivateKey(p).D.Bytes()
}

//PrivteKeyFromBytes makes a private key from a slice of bytes and reterns it.
func PrivteKeyFromBytes(b []byte) (priv *PrivateKey, err error) {
	if len(b) < 64 {
		return nil, fmt.Errorf("Could not parse commit bytes")
	}
	D := new(big.Int).SetBytes(b)
	priv = new(PrivateKey)
	priv.PublicKey.Curve = elliptic.P521()
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(b)
	priv.D = D
	return priv, nil
}

//PrivteKeyFromD makes a private key from a big int and reterns it
func PrivateKeyFromD(D big.Int) (*PrivateKey, error) {
	priv, err := PrivteKeyFromBytes(D.Bytes())
	return priv, err
	//priv := new(PrivateKey)
	//priv.PublicKey.Curve = elliptic.P521()
	//priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(D.Bytes())
	//priv.D = &D
	//return priv
}

//KeyGen makes a private new key and reterns it
func KeyGen() *PrivateKey {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	prikey := PrivateKey(*priv)
	return &prikey
}
