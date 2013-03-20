package main

//D *big.Int is all the key matereal save as bin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"
)

func KeyHash(p ecdsa.PrivateKey) []byte {
	var h hash.Hash = sha256.New()
	h.Write(KeyBytes(p))
	return h.Sum(nil)
}

func KeyBytes(p ecdsa.PrivateKey) []byte {
	return p.D.Bytes()
}

func getPiblicKeyForHkid(hkid []byte) *ecdsa.PublicKey {
	marshaledKey, _ := GetBlob(hkid)
	PublicKey := new(ecdsa.PublicKey)
	PublicKey.Curve = elliptic.P521()
	x, y := elliptic.Unmarshal(elliptic.P521(), marshaledKey)
	PublicKey.X, PublicKey.Y = x, y
	fmt.Print(PublicKey)
	return PublicKey
}

func getPrivateKeyForHkid(hkid []byte) (priv *ecdsa.PrivateKey, err error) {
	b, err := GetKey(hkid)
	priv = new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P521()
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(b)
	D := new(big.Int)
	priv.D = D.SetBytes(b)
	return
}

func PrivteKeyFromD(D big.Int) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P521()
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(D.Bytes())
	priv.D = &D
	return priv
}
