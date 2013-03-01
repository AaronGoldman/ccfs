package main

//D *big.Int is all the key matereal save as bin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

func getPiblicKeyForHkid(hkid [32]byte) *ecdsa.PublicKey {
	marshaledKey, _ := getBlob(hkid)
	PublicKey := new(ecdsa.PublicKey)
	PublicKey.Curve = elliptic.P521()
	x, y := elliptic.Unmarshal(elliptic.P521(), marshaledKey)
	PublicKey.X, PublicKey.Y = x, y
	fmt.Print(PublicKey)
	return PublicKey
}

func getPrivateKeyForHkid(hkid [32]byte) (priv *ecdsa.PrivateKey, err error) {
	b, err := getKey(hkid)
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
