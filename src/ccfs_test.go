package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestPath(t *testing.T) {
	//key for tag
	D := new(big.Int)
	D, _ = new(big.Int).SetString("3996811067068239799317867982525094232268699726722344370689730210711314983767775860556101498400185744208447673206609026128894016152514163591905578729891874833", 10)
	privT := PrivteKeyFromD(*D)
	keyT := elliptic.Marshal(privT.PublicKey.Curve,
		privT.PublicKey.X, privT.PublicKey.Y)
	hkidT, _ := hex.DecodeString(GenerateHKID(privT)) //gen HKID for tag
	PostKey(D.Bytes(), hkidT)                         //place key for tag	
	PostBlob(keyT)                                    //store public key

	//key for commit
	//r, ok := new(big.Int).SetString(s, 16)
	D, _ = new(big.Int).SetString("4629814823893296480016411334808793836186124559723200979962176753724976464088706463001383556112424820911870650421151988906751710824965155500230480521264034469", 10)
	privC := PrivteKeyFromD(*D)
	keyC := elliptic.Marshal(privC.PublicKey.Curve,
		privC.PublicKey.X, privC.PublicKey.Y)
	hkidC, _ := hex.DecodeString(GenerateHKID(privC)) //gen HKID for commit
	PostKey(D.Bytes(), hkidC)                         //place key for commit
	PostBlob(keyC)                                    //store public key

	//blob
	testBlob := []byte("testing") //gen test blob
	PostBlob(testBlob)            //store test blob

	//tag
	testTagPointingToTestBlob := NewTag(testBlob.Hash(),
		"blob",
		"testBlob",
		hkidT) //gen test tag
	PostTag(testTagPointingToTestBlob) //post test tag

	//list
	testListPiontingToTestTag := NewList(testTagPointingToTestBlob.Hash(),
		"tag",
		"testTag") //gen test list
	PostBlob(testListPiontingToTestTag) //store test list

	//commit
	testCommitPointingToTestList := NewCommit(testListPiontingToTestTag.Hash(),
		hkidC) //gen test commit
	PostCommit(testCommitPointingToTestList, version) //post test commit
}

/*func TestNewCommit(t *testing.T) {
	listHash := []byte{0xfa,
		0x84, 0xff, 0xaa, 0xe4, 0xd6, 0x5f, 0x49, 0x67, 0x67, 0x8c, 0x95, 0xb7, 0xf9, 0x6d, 0x61,
		0xe0, 0x67, 0x2b, 0x77, 0xe2, 0x67, 0x26, 0x78, 0x44, 0x44, 0x95, 0x24, 0x55, 0x07, 0x56, 0xb4}
	hkid := []byte{0xfa,
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
