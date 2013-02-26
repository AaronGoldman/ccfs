package main

//D *big.Int is all the key matereal save as bin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
)

func getPiblicKeyForHkid(hkid []byte) ecdsa.PublicKey {
	marshaledKey, _ := getBlob(hkid)
	PublicKey := new(ecdsa.PublicKey)
	PublicKey.Curve = elliptic.P521()
	x, y := elliptic.Unmarshal(elliptic.P521(), marshaledKey)
	PublicKey.X, PublicKey.Y = x, y
	return *PublicKey
}
