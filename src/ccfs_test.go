package main

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"bytes"
)

func TestPath(t *testing.T) {
	//key for tag
	D := new(big.Int)
	D, _ = new(big.Int).SetString("399681106706823979931786798252509423226869972"+
		"6722344370689730210711314983767775860556101498400185744208447673206609026"+
		"128894016152514163591905578729891874833", 10)
	privT := PrivteKeyFromD(*D)
	keyT := elliptic.Marshal(privT.PublicKey.Curve,
		privT.PublicKey.X, privT.PublicKey.Y)
	hkidT, _ := hex.DecodeString(GenerateHKID(privT)) //gen HKID for tag
	err := PostKey(privT)                             //place key for tag	
	if err != nil {
		fmt.Println(err)
	}
	err = PostBlob(keyT) //store public key
	if err != nil {
		fmt.Println(err)
	}

	//key for commit
	//r, ok := new(big.Int).SetString(s, 16)
	D, _ = new(big.Int).SetString("462981482389329648001641133480879383618612455"+
		"9723200979962176753724976464088706463001383556112424820911870650421151988"+
		"906751710824965155500230480521264034469", 10)
	privC := PrivteKeyFromD(*D)
	keyC := elliptic.Marshal(privC.PublicKey.Curve,
		privC.PublicKey.X, privC.PublicKey.Y)
	hkidC, _ := hex.DecodeString(GenerateHKID(privC)) //gen HKID for commit
	err = PostKey(privC)                              //place key for commit
	if err != nil {
		fmt.Println(err)
	}
	err = PostBlob(keyC) //store public key
	if err != nil {
		fmt.Println(err)
	}

	//Post blob
	testBlob := blob([]byte("testing")) //gen test blob
	err = PostBlob(testBlob)            //store test blob
	if err != nil {
		fmt.Println(err)
	}

	//post tag
	testTagPointingToTestBlob := NewTag(testBlob.Hash(),
		"blob",
		"testBlob",
		hkidT) //gen test tag
	err = PostTag(testTagPointingToTestBlob) //post test tag
	if err != nil {
		fmt.Println(err)
	}

	//post list
	testListPiontingToTestTag := NewList(testTagPointingToTestBlob.Hkid(),
		"tag",
		"testTag") //gen test list
	err = PostBlob(testListPiontingToTestTag.Bytes()) //store test list
	if err != nil {
		fmt.Println(err)
	}

	// post commit
	testCommitPointingToTestList := NewCommit(testListPiontingToTestTag.Hash(),
		hkidC) //gen test commit
	err = PostCommit(testCommitPointingToTestList) //post test commit
	if err != nil {
		fmt.Println(err)
	}

	//print
	fmt.Printf("%s\n%s\n%s\n%s\n",
		hex.EncodeToString(testBlob.Hash()),
		hex.EncodeToString(testTagPointingToTestBlob.Hash()),
		hex.EncodeToString(testListPiontingToTestTag.Hash()),
		hex.EncodeToString(testCommitPointingToTestList.Hash()))

	//get commit
	hkid, _ := hex.DecodeString("1312ac161875b270da2ae4e1471ba94a" +
		"9883419250caa4c2f1fd80a91b37907e")
	commitbytes, err := GetCommit(hkid)
	if err != nil {
		panic(err)
	}
	commitStrings := strings.Split(string(commitbytes), ",\n")
	listHash, _ := hex.DecodeString(commitStrings[0])
	version, _ := strconv.ParseInt(commitStrings[1], 10, 64)
	chkid, _ := hex.DecodeString(commitStrings[2])
	signature, _ := hex.DecodeString(commitStrings[3])
	testcommit := commit{listHash, version, chkid, signature}
	//fmt.Println(testcommit)
	fmt.Printf("authentic commit:%v\n", testcommit.Verifiy())
	//get list
	listbytes, err := GetBlob(listHash)
	if err != nil {
		panic(err)
	}
	listEntries := strings.Split(string(listbytes), "\n")
	entries := []entry{}
	cols := []string{}
	for _, element := range listEntries {
		cols = strings.Split(element, ",")
		entryHash, _ := hex.DecodeString(cols[0])
		entryTypeString := cols[1]
		entryNameSegment := cols[2]
		entries = append(entries, entry{entryHash, entryTypeString, entryNameSegment})
	}
	testlist := list(entries)
	fmt.Printf("authentic list:%v\n", bytes.Equal(listHash,testlist.Hash()))
	//get tag
	//tagBytes := GetTag(entryHash,"testBlob")
	
	//get blob
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

/*
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
*/

/*
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
*/
/*
func TestTaggen(t *testing.T) {
	priv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//test(priv)
	fmt.Printf("\nTestTaggen\n%v\n", GenerateTag([]byte("testing"), "blob", "test", priv))
}
*/
/*
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
*/
