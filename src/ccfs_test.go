package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
)

func BenchmarkStoreOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ioutil.WriteFile("../storeone", []byte("storeone"), 0664)
		if err != nil {
			panic(err)
		}
		_, err = ioutil.ReadFile("../storeone")
		if err != nil {
			panic(err)
		}
		//fmt.Print(data)
	}
}

func BenchmarkPath(b *testing.B) {

	//key for tag
	D := new(big.Int)
	D, _ = new(big.Int).SetString("399681106706823979931786798252509423226869972"+
		"6722344370689730210711314983767775860556101498400185744208447673206609026"+
		"128894016152514163591905578729891874833", 10) //this number is the key
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

	for i := 0; i < b.N; i++ {
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
		//fmt.Printf("%s\n%s\n%s\n%s\n",
		//	hex.EncodeToString(testBlob.Hash()),
		//	hex.EncodeToString(testTagPointingToTestBlob.Hash()),
		//	hex.EncodeToString(testListPiontingToTestTag.Hash()),
		//	hex.EncodeToString(testCommitPointingToTestList.Hash()))

		//get commit
		hkid, _ := hex.DecodeString("1312ac161875b270da2ae4e1471ba94a" +
			"9883419250caa4c2f1fd80a91b37907e")
		testcommit, err := GetCommit(hkid)
		if err != nil {
			panic(err)
		}

		fmt.Printf("authentic commit:%v\n", testcommit.Verifiy())

		//get list
		listbytes, err := GetBlob(testcommit.listHash)
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
		fmt.Printf("authentic list:%v\n", bytes.Equal(testcommit.listHash, testlist.Hash()))
		//get tag
		testTag, err := GetTag(testlist[0].Hash, "testBlob")
		fmt.Printf("authentic tag:%v\n", testTag.Verifiy())
		//get blob
		testBlob, _ = GetBlob(testTag.HashBytes)
		fmt.Printf("authentic blob:%v\n", bytes.Equal(testTag.HashBytes, testBlob.Hash()))
	}
}

func TestKeyGen(b *testing.T) {
	c := elliptic.P521()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		b.Errorf("Error %v", err)
	}
	fmt.Printf("TestKeyGen\nX = %v\nY = %v\nD = %v\n", priv.PublicKey.X, priv.PublicKey.Y, priv.D)
	err = PostKey(priv)
	if err != nil {
		b.Errorf("Error %v", err)
	}
	PostBlob(elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y))
}
