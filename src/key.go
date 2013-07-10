package main

//D *big.Int is all the key matereal save as bin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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
	curve := elliptic.P521()
	x, y := elliptic.Unmarshal(elliptic.P521(), marshaledKey)
	publicKey := ecdsa.PublicKey{
		curve, //elliptic.Curve
		x,     //big.Int
		y}     //big.Int
	return &publicKey
}

func getPrivateKeyForHkid(hkid HKID) (priv *ecdsa.PrivateKey, err error) {
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

func HkidFromD(D big.Int) HKID {
	priv := PrivteKeyFromD(D)
	key := elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y)
	hkid := GenerateHKID(priv)
	err := PostKey(priv) //store privet key
	if err != nil {
		panic(err)
	}
	err = PostBlob(key) //store public key
	if err != nil {
		panic(err)
	}
	return hkid
}

func hkidFromDString(str string, base int) HKID {
	D, success := new(big.Int).SetString(str, base)
	if !success {
		panic(nil)
	}
	return HkidFromD(*D)
}

func KeyGen() *ecdsa.PrivateKey {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return priv
}
