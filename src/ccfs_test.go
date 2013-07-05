package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGet(t *testing.T) {
	dabytes, err := hex.DecodeString("1312ac161875b270da2ae4e1471ba94a" +
		"9883419250caa4c2f1fd80a91b37907e")
	hkid := HKID(dabytes)
	path := "testTag/testBlob"
	b := []byte(":(")
	if err == nil {
		fmt.Println(hkid.Hex())
		b, err = Get(hkid, path)
	}
	if !bytes.Equal([]byte("testing"), b) || err != nil {
		t.Fail()
	}
}

/*//test for new post flow
func TestPath(t *testing.T) {
	repoHkid := hkidFromDString("399681106706823979931786798252509423226869972672"+
		"2344370689730210711314983767775860556101498400185744208447673206609026128"+
		"894016152514163591905578729891874833", 10)
	domainHkid := hkidFromDString("462981482389329648001641133480879383618612455972"+
		"3200979962176753724976464088706463001383556112424820911870650421151988906"+
		"751710824965155500230480521264034469", 10)
	err := InitRepo(repoHkid) // post commit //post list
	if err != nil {
		panic(err)
	}

	err = InitDomain(domainHkid) // post tag
	if err != nil {
		panic(err)
	}
	InsertDomain(repoHkid, domainHkid, "testTag")
	// Post blob
	err = Post(repoHkid, "testTag/testBlob", blob([]byte("testing")))
	if err != nil {
		panic(err)
	}
} //*/

func BenchmarkStoreOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ioutil.WriteFile("../bin/storeone", []byte("storeone"), 0664)
		if err != nil {
			panic(err)
		}
		_, err = ioutil.ReadFile("../bin/storeone")
		if err != nil {
			panic(err)
		}
		//fmt.Print(data)
	}
}

func BenchmarkPath(b *testing.B) {
	//key for tag
	hkidT := hkidFromDString("39968110670682397993178679825250942322686997267223"+
		"4437068973021071131498376777586055610149840018574420844767320660902612889"+
		"4016152514163591905578729891874833", 10)

	//key for commit
hkidC := hkidFromDString("4629814823893296480016411334808793836186124559723200"+
		"9799621767537249764640887064630013835561124248209118706504211519889067517"+
		"10824965155500230480521264034469", 10)

	for i := 0; i < b.N; i++ {
		//Post blob
		testBlob := blob([]byte("testing")) //gen test blob
		err := PostBlob(testBlob)            //store test blob
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
		//commit -> list -> tag -> blob

		//get commit
		hkid, _ := hex.DecodeString("1312ac161875b270da2ae4e1471ba94a" +
			"9883419250caa4c2f1fd80a91b37907e")
		testcommit, err := GetCommit(hkid)
		if err != nil {
			panic(err)
		}

		//fmt.Printf("authentic commit:%v\n", testcommit.Verifiy())
		if !testcommit.Verifiy() {
			b.FailNow()
		}
		//get list
		listbytes, err := GetBlob(testcommit.listHash)
		if err != nil {
			panic(err)
		}

		testlist := NewListFromBytes(listbytes)

		//fmt.Printf("authentic list:%v\n", )
		if !bytes.Equal(testcommit.listHash, testlist.Hash()) {
			b.FailNow()
		}
		//get tag
		testTag, err := GetTag(testlist[0].Hash, "testBlob")
		//fmt.Printf("authentic tag:%v\n", testTag.Verifiy())
		if !testTag.Verifiy() {
			b.FailNow()
		}
		//get blob
		testBlob, _ = GetBlob(testTag.HashBytes)
		//fmt.Printf("authentic blob:%v\n", bytes.Equal(testTag.HashBytes,
		//	testBlob.Hash()))
		if !bytes.Equal(testTag.HashBytes, testBlob.Hash()) {
			b.FailNow()
		}
	}
}

func TestKeyGen(b *testing.T) {
	c := elliptic.P521()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		b.Errorf("Error %v", err)
	}
	//fmt.Printf("TestKeyGen\nX = %v\nY = %v\nD = %v\n", priv.PublicKey.X,
	//	priv.PublicKey.Y, priv.D)
	err = PostKey(priv)
	if err != nil {
		b.Errorf("Error %v", err)
	}
	PostBlob(elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y))
}
