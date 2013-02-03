package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io"
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
