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
	testpath := "TestPostBlob"
	indata := []byte("TestPostData")
	Post(testhkid, testpath, BlobFromBytes(indata))
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
	//log.Printf("\n\tkey: %s\n\tpath: %s\n\tindata: %s\n\toutdata: %s\n",
	//	testhkid.Hex(),
	//	testpath,
	//	indata,
	//	outdata)
}

func TestPostListListBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
	testpath := "TestPostList1/TestPostList2/TestPostBlob"
	indata := []byte("TestPostListListBlobData")
	Post(testhkid, testpath, BlobFromBytes(indata))
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostCommitBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("65232373562705602286177837897283"+
		"2941659551264911224937349783059207224141689361121606942380473043"+
		"7860475300564272976762085068519188612732562106886379081213385", 10)
	testRepoHkid := hkidFromDString("22371143209450593169269383669277"+
		"1410459232098247632372342448006863927240156318431751613873181811"+
		"3842641285036340692759591635625837820111726090732747634977413", 10)
	InsertRepo(testhkid, "TestPostCommit", testRepoHkid)
	testpath := "TestPostCommit/TestPostBlob"
	indata := []byte("TestPostCommitBlobData")
	Post(testhkid, testpath, BlobFromBytes(indata))
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostTagBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("65232373562705602286177837897283"+
		"2941659551264911224937349783059207224141689361121606942380473043"+
		"7860475300564272976762085068519188612732562106886379081213385", 10)
	testDomainHkid := hkidFromDString("48133155371863219707191657794884"+
		"7537768863308478273117048217437472569476796507874261675750733630"+
		"0675115007319549300204862762937376227751462258339344895829332", 10)
	InsertDomain(testhkid, "TestPostTag", testDomainHkid)
	testpath := "TestPostTag/TestPostBlob"
	indata := []byte("TestPostTagBlobData")
	Post(testhkid, testpath, BlobFromBytes(indata))
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostTagCommitBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domainHkid := hkidFromDString("39968110670682397993178679825250"+
		"9423226869972672234437068973021071131498376777586055610149840018"+
		"5744208447673206609026128894016152514163591905578729891874833", 10)
	repoHkid := hkidFromDString("94522678075182002377842140271746"+
		"1502019766737301062946423046280258817349516439546479625226895211"+
		"80808353215150536034481206091147220911087792299373183736254", 10)
	err := InsertDomain(testhkid, "testTag", domainHkid)
	err = InsertRepo(testhkid, "testTag/testCommit", repoHkid)
	indata := blob([]byte("TestTagCommitBlobData"))
	testpath := "testTag/testCommit/testBlob"
	_, err = Post(testhkid, testpath, indata)
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostTagTagBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domain1Hkid := hkidFromDString("32076859881811206392323279987831"+
		"3732334949433938278619036381396945204532895319697233476900342324"+
		"1692155627674142765782672165943417419038514237233188152538761", 10)
	domain2Hkid := hkidFromDString("49220288257701056900010210884990"+
		"0714973444364727181180850528073586453638681999434006549298762097"+
		"4197649255374796716934112121800838847071661501215957753532505", 10)
	err := InsertDomain(testhkid, "testTag1", domain1Hkid)
	err = InsertDomain(testhkid, "testTag1/testTag2", domain2Hkid)
	indata := blob([]byte("TestTagTagBlobData"))
	testpath := "testTag1/testTag2/testBlob"
	_, err = Post(testhkid, testpath, indata)
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostListCommitBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	repoHkid := hkidFromDString("59803288656043806807393139191118"+
		"0777289273091938777159029927847596771500408478956278378281366717"+
		"7487960901581583946753338859223459810645621124266443931192097", 10)
	err := InsertDomain(testhkid, "testList/testCommit", repoHkid)
	indata := blob([]byte("TestListCommitBlobData"))
	testpath := "testList/testCommit/testBlob"
	_, err = Post(testhkid, testpath, indata)
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostTagListBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domainHkid := hkidFromDString("62089221762245310704629142682144"+
		"1944826557905230450143203631438168806532495876980559885034903315"+
		"4294997505754401230560960060918213268981906409591978967796584", 10)
	err := InsertDomain(testhkid, "testTag", domainHkid)
	indata := blob([]byte("TestTagListBlobData"))
	testpath := "testTag/testList/testBlob"
	_, err = Post(testhkid, testpath, indata)
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostListTagBlob(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domainHkid := hkidFromDString("12796633883654089746486670711069"+
		"9781359720828332046318301886846633714179790444071153863702142701"+
		"146332294245448463914286494124849121460550667767568731696934", 10)
	err := InsertDomain(testhkid, "testList/testTag", domainHkid)
	indata := blob([]byte("TestListTagBlobData"))
	testpath := "testList/testTag/testBlob"
	_, err = Post(testhkid, testpath, indata)
	outdata, err := Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
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

func TestKeyGen(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	t.SkipNow()
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
