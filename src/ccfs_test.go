package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

/*
func TestPath(t *testing.T) {
	//key
	postKey(keyT)//place key for tag
	hkidT = keyT.Hash()//gen HKID for tag
	postKey(keyC)//place key for commit
	hkidC = keyC.Hash()//gen HKID for commit

	//blob
	testBlob := newBlob([]byte("testing"))//gen test blob
	postBlob(testBlob)//store test blob

	//tag
	testTagPointingToTestBlob := newTag(testBlob.Hash(),
		"blob",
		"testBlob",
		tagVersion,
		hkidT)//gen test tag
	postTag(testTagPointingToTestBlob)//post test tag

	//list
	testListPiontingToTestTag = NewList(testTagPointingToTestBlob.Hash(),
		"tag",
		"testTag")//gen test list
	postBlob(testListPiontingToTestTag)//store test list

	//commit
	testCommitPointingToTestList =: NewCommit(testListPiontingToTestTag.Hash(),
	 hkidC)//gen test commit
	postCommit(testCommitPointingToTestList, version)//post test commit
}
*/

/*func TestNewCommit(t *testing.T) {
	listHash := [32]byte{0xfa,
		0x84, 0xff, 0xaa, 0xe4, 0xd6, 0x5f, 0x49, 0x67, 0x67, 0x8c, 0x95, 0xb7, 0xf9, 0x6d, 0x61,
		0xe0, 0x67, 0x2b, 0x77, 0xe2, 0x67, 0x26, 0x78, 0x44, 0x44, 0x95, 0x24, 0x55, 0x07, 0x56, 0xb4}
	hkid := [32]byte{0xfa,
		0x84, 0xff, 0xaa, 0xe4, 0xd6, 0x5f, 0x49, 0x67, 0x67, 0x8c, 0x95, 0xb7, 0xf9, 0x6d, 0x61,
		0xe0, 0x67, 0x2b, 0x77, 0xe2, 0x67, 0x26, 0x78, 0x44, 0x44, 0x95, 0x24, 0x55, 0x07, 0x56, 0xb4}
	fmt.Print(listHash, hkid)
	c := NewCommit(listHash, hkid)
	fmt.Print(c.Verifiy())
}*/

func TestKeyGen(t *testing.T) {
	c := elliptic.P521()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	fmt.Printf("TestKeyGen\nX = %v\nY = %v\nD = %v\n", priv.PublicKey.X, priv.PublicKey.Y, priv.D)
	//priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(priv.D.Bytes())
	//fmt.Printf("X = %v\nY = %v\nD = %v\n", priv.PublicKey.X, priv.PublicKey.Y, priv.D)
	test(priv, t)
	//SaveKey(priv)
	SavePrivateKey(priv)
	priv, err = LoadPrivateKey(GenerateHKID(priv))
	fmt.Printf("\nhkid:%v\n", GenerateHKID(priv))
}

func TestCommitgen(t *testing.T) {
	priv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	SavePrivateKey(priv)
	//fmt.Printf("\nTestCommitgen\nX = %v\nY = %v\nD = %v\n", priv.PublicKey.X, priv.PublicKey.Y, priv.D)
	hkid := GenerateHKID(priv)
	reconstructedkey, err := LoadPrivateKey(hkid)
	if err != nil {
		panic(err)
	}
	reconstructedhkid := GenerateHKID(reconstructedkey)
	fmt.Printf("\n%v\n%v", hkid, reconstructedhkid)
	fmt.Printf("\nTestCommitgen\n%v", GenerateCommit([]byte("testing"), priv))
}

func TestTaggen(t *testing.T) {
	priv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//test(priv)
	fmt.Printf("\nTestTaggen\n%v\n", GenerateTag([]byte("testing"), "blob", "test", priv))
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
