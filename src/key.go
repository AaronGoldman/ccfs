package main

//D *big.Int is all the key matereal save as bin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

func getPiblicKeyForHkid(hkid []byte) ecdsa.PublicKey {
	marshaledKey, _ := getBlob(hkid)
	PublicKey := new(ecdsa.PublicKey)
	PublicKey.Curve = elliptic.P521()
	x, y := elliptic.Unmarshal(elliptic.P521(), marshaledKey)
	PublicKey.X, PublicKey.Y = x, y
	return *PublicKey
}

func PrivteKeyFromD(D big.Int) ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P521()
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(D.Bytes())
	priv.D = &D
	return *priv
}
