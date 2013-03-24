package main

//import (
//	"crypto/ecdsa"
//	"crypto/elliptic"
//	"encoding/pem"
//	"fmt"
//	"io/ioutil"
//	"math/big"
//)

//type Privkey ecdsa.PrivateKey

//type Pubkey ecdsa.PublicKey

//func SavePrivateKey(priv *ecdsa.PrivateKey) (err error) {
//	fmt.Print(GenerateHKID(priv))
//	filepath := fmt.Sprintf("../keys/%s", GenerateHKID(priv))
//	Block := pem.Block{Type: "EC PRIVATE KEY", Bytes: priv.D.Bytes()}
//	err = ioutil.WriteFile(filepath, pem.EncodeToMemory(&Block), 0440)
//	return
//}

//func LoadPrivateKey(hkid string) (priv *ecdsa.PrivateKey, err error) {
//	filepath := fmt.Sprintf("../keys/%s", hkid)
//	data, err := ioutil.ReadFile(filepath)
//	if err != nil {
//		panic(err)
//	}
//	b, _ := pem.Decode(data)
//	priv = new(ecdsa.PrivateKey)
//	priv.PublicKey.Curve = elliptic.P521()
//	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P521().ScalarBaseMult(b.Bytes)
//	D := new(big.Int)
//	priv.D = D.SetBytes(b.Bytes)
//	return
//}
