package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestKeyGen(t *testing.T) {
	c := elliptic.P521()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		t.Errorf("Error", err)
	}
	fmt.Printf("X = %v\nY = %v\nD = %v\n", priv.PublicKey.X, priv.PublicKey.Y, priv.D)
	//priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(priv.D.Bytes())
	//fmt.Printf("X = %v\nY = %v\nD = %v\n", priv.PublicKey.X, priv.PublicKey.Y, priv.D)
	test(priv, t)
	//SaveKey(priv)
	SavePrivateKey(priv)
	priv, err = LoadPrivateKey(GenerateHKID(priv))
	fmt.Printf("\nhkid:%v\n", GenerateHKID(priv))
}

//func Commitgen_test(t *testing.T) {
//priv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
//SaveKey(priv)
//fmt.Print("X = %v\nY = %v\nD = %v\n",priv.PublicKey.X,priv.PublicKey.Y,priv.D)
//hkid := GenerateHKID(priv)
//reconstructedkey := LoadKey(hkid)
//reconstructedhkid := GenerateHKID(reconstructedkey)
//fmt.Printf("\n%v\n%v", hkid, reconstructedhkid)
//fmt.Printf("%v", GenerateCommit([]byte("testing"), priv))
//}

func Taggen_test(t *testing.T) {
	priv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//test(priv)
	fmt.Printf("%v\n", GenerateTag([]byte("testing"), "blob", "test", priv))
}

func test(priv *ecdsa.PrivateKey, t *testing.T) {
	mar := elliptic.Marshal(
		priv.PublicKey.Curve,
		priv.PublicKey.X,
		priv.PublicKey.Y)
	x, y := elliptic.Unmarshal(elliptic.P521(), mar)
	maredPublicKey := new(ecdsa.PublicKey)
	maredPublicKey.Curve = elliptic.P521()
	maredPublicKey.X = x
	maredPublicKey.Y = y

	hashed := []byte("testing")
	r, s, _ := ecdsa.Sign(rand.Reader, priv, hashed)
	valid := ecdsa.Verify(maredPublicKey, hashed, r, s)
	if valid != true {
		t.Errorf("failed Verify\n")
	}
	invalid := ecdsa.Verify(maredPublicKey, []byte("fail testing"), r, s)
	if invalid != false {
		t.Errorf("failed falsify\n")
	}
	//fmt.Printf("valid:%v in:%v marsize:%v bits\n\n\n", valid, invalid, len(mar)*8)
}
