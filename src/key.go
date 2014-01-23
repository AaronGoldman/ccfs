package main

//D *big.Int is all the key matereal save as bin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"log"
	"math/big"
)

type PrivateKey ecdsa.PrivateKey

type PublicKey ecdsa.PublicKey

//Bytes returns the marshaled public key as a slice of byte.
func (p PublicKey) Bytes() []byte {
	return elliptic.Marshal(p.Curve, p.X, p.Y)
}

//Hkid returns the hkid for the public key.
func (p PublicKey) Hkid() HKID {
	var h hash.Hash = sha256.New()
	h.Write(p.Bytes())
	return h.Sum(make([]byte, 0))
}

//Hkid returns the hkid that for private key.
//this is the hcid of your public key
func (p PrivateKey) Hkid() HKID {
	var h hash.Hash = sha256.New()
	h.Write(elliptic.Marshal(p.PublicKey.Curve, p.PublicKey.X, p.PublicKey.Y))
	return h.Sum(make([]byte, 0))
}

//Hash returns the hcid for the PrivateKey
func (p PrivateKey) Hash() []byte {
	var h hash.Hash = sha256.New()
	h.Write(p.Bytes())
	return h.Sum(nil)
}

//returns true if PrivateKey and PrivateKey..PublicKey are a pair.
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

//getPiblicKeyForHkid uses the lookup services to get a public key for an hkid
func getPiblicKeyForHkid(hkid HKID) PublicKey {
	marshaledKey, err := GetBlob(HCID(hkid))
	if err != nil {
		return PublicKey{}
	}
	curve := elliptic.P521()
	x, y := elliptic.Unmarshal(elliptic.P521(), marshaledKey)
	pubKey := ecdsa.PublicKey{
		curve, //elliptic.Curve
		x,     //X *big.Int
		y}     //Y *big.Int
	return PublicKey(pubKey)
}

//getPrivateKeyForHkid uses the lookup services to get a private key for an hkid
func getPrivateKeyForHkid(hkid HKID) (k *PrivateKey, err error) {
	k, err = GetKey(hkid)
	return k, err
}

//PrivteKeyFromBytes makes a private key from a slice of bytes and reterns it.
func PrivteKeyFromBytes(b []byte) *PrivateKey {
	D := new(big.Int).SetBytes(b)
	priv := new(PrivateKey)
	priv.PublicKey.Curve = elliptic.P521()
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(b)
	priv.D = D
	return priv
}

func PrivteKeyFromD(D big.Int) *PrivateKey {
	return PrivteKeyFromBytes(D.Bytes())
	//priv := new(PrivateKey)
	//priv.PublicKey.Curve = elliptic.P521()
	//priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(D.Bytes())
	//priv.D = &D
	//return priv
}

func HkidFromD(D big.Int) HKID {
	priv := PrivteKeyFromD(D)
	key := elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y)
	hkid := priv.Hkid()
	err := PostKey(priv) //store privet key
	if err != nil {
		log.Panic(err)
	}
	err = PostBlob(key) //store public key
	if err != nil {
		log.Panic(err)
	}
	return hkid
}

func hkidFromDString(str string, base int) HKID {
	D, success := new(big.Int).SetString(str, base)
	if !success {
		log.Panic(nil)
	}
	return HkidFromD(*D)
}

func KeyGen() *PrivateKey {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	prikey := PrivateKey(*priv)
	return &prikey
}

func GenHKID() HKID {
	privkey := KeyGen()
	PostBlob(elliptic.Marshal(privkey.PublicKey.Curve,
		privkey.PublicKey.X, privkey.PublicKey.Y))
	return privkey.Hkid()
}
