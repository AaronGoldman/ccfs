package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

type Privkey struct {
	Pubkey
	D *big.Int
}

type Pubkey struct {
	X, Y *big.Int
}

func SavePrivateKey(priv *ecdsa.PrivateKey) (err error) {
	fmt.Print(GenerateHKID(priv))
	filepath := fmt.Sprintf("../keys/%s", GenerateHKID(priv))
	fo, err := os.Create(filepath)
	if err != nil {
		return
	}
	Block := pem.Block{Type: "EC PRIVATE KEY", Bytes: priv.D.Bytes()}
	err = pem.Encode(fo, &Block)
	return
}

func LoadPrivateKey(hkid string) (priv *ecdsa.PrivateKey, err error) {
	filepath := fmt.Sprintf("../keys/%s", hkid)
	//fmt.Print(filepath)
	//fi, err := os.Open(filepath)
	//if err != nil {
	//	return
	//	}
	//defer fi.Close()
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	b, _ := pem.Decode(data)
	//fmt.Print("\nData:",data,"\nB:",b.Bytes,"\nRest:",rest,"\n")
	priv = new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P521()
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(b.Bytes)
	D := new(big.Int)
	priv.D = D.SetBytes(b.Bytes)
	//fmt.Print(priv)
	return
}

func SaveKey(priv *ecdsa.PrivateKey) {
	filepath := fmt.Sprintf("../keys/%s", GenerateHKID(priv))
	//fmt.Print(*priv)
	fo, err := os.Create(filepath)

	if err != nil {
		//fmt.Printf("%v", err)
		//return
		panic(err)
	}
	enc := gob.NewEncoder(fo)
	pu := Pubkey{priv.PublicKey.X, priv.PublicKey.Y}
	pr := Privkey{pu, priv.D}
	err = enc.Encode(pr)
	if err != nil {
		fmt.Printf("\n%v", err)
	}
	fo.Close()
}

func LoadKey(hkid string) (priv *ecdsa.PrivateKey) {
	filepath := fmt.Sprintf("../keys/%s", hkid)
	fmt.Print(filepath)
	fi, err := os.Create(filepath)

	if err != nil {
		fmt.Printf("\n%v", err)
		return
		//panic(err)
	}
	dec := gob.NewDecoder(fi)
	var prikey Privkey
	err = dec.Decode(&prikey)
	fi.Close()
	if err != nil && err != io.EOF {
		panic(fmt.Sprintf("%v", err))
	}
	priv = &ecdsa.PrivateKey{ecdsa.PublicKey{elliptic.P521(), prikey.Pubkey.X, prikey.Pubkey.Y}, prikey.D}
	fmt.Printf("%v\n", priv)
	return priv
}
