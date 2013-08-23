package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestPostBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
	indata := []byte("TestPostBlob")
	log.Printf("\n\tkey: %s\n", testhkid.Hex())
	Post(testhkid, "TestPostBlob", BlobFromBytes(indata))
	log.Print("Posted")
	outdata, err := Get(testhkid, "TestPostBlob")
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
	log.Print("Got")
}

func TestGet(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	//hkidFromDString("4629814823893296480016411334808793836186124559723200"+
	//	"9799621767537249764640887064630013835561124248209118706504211519889067517"+
	//	"10824965155500230480521264034469", 10) //store the key
	setup_for_gets()

	dabytes, err := hex.DecodeString("1312ac161875b270da2ae4e1471ba94a" +
		"9883419250caa4c2f1fd80a91b37907e")

	hkid := HKID(dabytes)
	path := "testTag/testBlob"
	b := []byte(":(")
	if err == nil {
		//log.Println(hkid.Hex())
		b, err = Get(hkid, path)
	}
	if !bytes.Equal([]byte("testing"), b) || err != nil {
		t.Fail()
	}
}

//test for new post flow
func DontTestPost(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	//t.SkipNow()
	repoHkid := hkidFromDString("46298148238932964800164113348087938361861245597"+
		"2320097996217675372497646408870646300138355611242482091187065042115198890"+
		"6751710824965155500230480521264034469", 10)
	domainHkid := hkidFromDString("399681106706823979931786798252509423226869972"+
		"6722344370689730210711314983767775860556101498400185744208447673206609026"+
		"128894016152514163591905578729891874833", 10)
	err := InitRepo(repoHkid) // post commit //post list
	if err != nil {
		log.Panic(err)
	}
	log.Println("InitRepo")

	err = InitDomain(domainHkid, "null") // post tag
	if err != nil {
		log.Panic(err)
	}
	log.Println("InitDomain")
	InsertDomain(repoHkid, domainHkid, "testTag")
	log.Println("InsertDomain")
	// Post blob
	_, err = Post(repoHkid, "testTag/testBlob2", blob([]byte("testing2")))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Posted")
} //*/

func BenchmarkStoreOne(b *testing.B) {
	log.SetFlags(log.Lshortfile)
	for i := 0; i < b.N; i++ {
		err := os.MkdirAll("../bin/", 0764)
		err = ioutil.WriteFile("../bin/storeone", []byte("storeone"), 0664)
		if err != nil {
			log.Panic(err)
		}
		_, err = ioutil.ReadFile("../bin/storeone")
		if err != nil {
			log.Panic(err)
		}
		//log.Print(data)
	}
}

func BenchmarkPath(b *testing.B) {
	log.SetFlags(log.Lshortfile)
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
		err := PostBlob(testBlob)           //store test blob
		if err != nil {
			log.Println(err)
		}

		//post tag
		testTagPointingToTestBlob := NewTag(HID(testBlob.Hash()),
			"blob",
			"testBlob",
			hkidT) //gen test tag
		err = PostTag(testTagPointingToTestBlob) //post test tag
		if err != nil {
			log.Println(err)
		}

		//post list
		testListPiontingToTestTag := NewList(testTagPointingToTestBlob.Hkid(),
			"tag",
			"testTag") //gen test list
		err = PostBlob(testListPiontingToTestTag.Bytes()) //store test list
		if err != nil {
			log.Println(err)
		}

		// post commit
		testCommitPointingToTestList := NewCommit(testListPiontingToTestTag.Hash(),
			hkidC) //gen test commit
		err = PostCommit(testCommitPointingToTestList) //post test commit
		if err != nil {
			log.Println(err)
		}

		//print
		//log.Printf("%s\n%s\n%s\n%s\n",
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
			log.Panic(err)
		}

		//log.Printf("authentic commit:%v\n", testcommit.Verifiy())
		if !testcommit.Verifiy() {
			b.FailNow()
		}
		//get list
		listbytes, err := GetBlob(testcommit.listHash)
		if err != nil {
			log.Panic(err)
		}

		testlist := NewListFromBytes(listbytes)

		//log.Printf("authentic list:%v\n", )
		if !bytes.Equal(testcommit.listHash, testlist.Hash()) {
			b.FailNow()
		}
		//get tag
		_, testTagHash := testlist.hash_for_namesegment("testTag")
		testTag, err := GetTag(testTagHash.Bytes(), "testBlob")
		//log.Printf("authentic tag:%v\n", testTag.Verifiy())
		if !testTag.Verifiy() {
			b.FailNow()
		}
		//get blob
		testBlob, _ = GetBlob(HCID(testTag.HashBytes))
		//log.Printf("authentic blob:%v\n", bytes.Equal(testTag.HashBytes,
		//	testBlob.Hash()))
		if !bytes.Equal(testTag.HashBytes, testBlob.Hash()) {
			b.FailNow()
		}
	}
}

func DontTestKeyGen(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	//t.SkipNow
	c := elliptic.P521()
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	log.Printf("TestKeyGen\nX = %v\nY = %v\nD = %v\n", priv.PublicKey.X,
		priv.PublicKey.Y, priv.D)
	err = PostKey(priv)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	PostBlob(elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y))
}

func setup_for_gets() {
	hkidT := hkidFromDString("39968110670682397993178679825250942322686997267223"+
		"4437068973021071131498376777586055610149840018574420844767320660902612889"+
		"4016152514163591905578729891874833", 10)

	//key for commit
	hkidC := hkidFromDString("4629814823893296480016411334808793836186124559723200"+
		"9799621767537249764640887064630013835561124248209118706504211519889067517"+
		"10824965155500230480521264034469", 10)

	//Post blob
	testBlob := blob([]byte("testing")) //gen test blob
	err := PostBlob(testBlob)           //store test blob
	if err != nil {
		log.Println(err)
	}

	//post tag
	testTagPointingToTestBlob := NewTag(HID(testBlob.Hash()),
		"blob",
		"testBlob",
		hkidT) //gen test tag
	err = PostTag(testTagPointingToTestBlob) //post test tag
	if err != nil {
		log.Println(err)
	}

	//post list
	testListPiontingToTestTag := NewList(testTagPointingToTestBlob.Hkid(),
		"tag",
		"testTag") //gen test list
	err = PostBlob(testListPiontingToTestTag.Bytes()) //store test list
	if err != nil {
		log.Println(err)
	}

	// post commit
	testCommitPointingToTestList := NewCommit(testListPiontingToTestTag.Hash(),
		hkidC) //gen test commit
	err = PostCommit(testCommitPointingToTestList) //post test commit
	if err != nil {
		log.Println(err)
	}
}
