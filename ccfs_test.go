package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/hex"
	"log"
	"testing"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/localfile"
	"github.com/AaronGoldman/ccfs/services/timeout"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	localfile.Start()
	timeout.Start()
	objects.RegisterGeterPoster(
		services.GetPublicKeyForHkid,
		services.GetPrivateKeyForHkid,
		services.PostKey,
		services.PostBlob,
	)
}

func TestLowLevel(t *testing.T) {
	indata := objects.Blob([]byte("TestPostData"))
	services.PostBlob(indata)
	blobhcid, err := objects.HcidFromHex(
		"ca4c4244cee2bd8b8a35feddcd0ba36d775d68637b7f0b4d2558728d0752a2a2",
	)
	b, err := services.GetBlob(blobhcid)
	if !bytes.Equal(b.Hash(), blobhcid) {
		t.Fatalf("GetBlob Fail\nExpected: %s\nGot: %s\n", blobhcid, b.Hash()) //Changed
	}

	//6dedf7e580671bd90bc9d1f735c75a4f3692b697f8979a147e8edd64fab56e85
	testhkid := objects.HkidFromDString(
		"6523237356270560228617783789728329416595512649112249373497830592"+
			"0722414168936112160694238047304378604753005642729767620850685191"+
			"88612732562106886379081213385", 10)
	testCommit := objects.NewCommit(b.Hash(), testhkid)
	services.PostCommit(testCommit)
	c, err := services.GetCommit(testhkid)
	if err != nil {
		t.Fatalf("Get Commit Fail: %s", err)
	} else if !bytes.Equal(c.Hkid, testhkid) {
		t.Fatalf("Expected: %s Got: %s", testhkid, c.Hkid)
	} else if !c.Verify() {
		t.Fatalf("Commit Signature Invalid.")
	}

	//ede7bec713c93929751f18b1db46d4be3c95286bd5f2d92b9759ff02115dc312
	taghkid := objects.HkidFromDString(
		"4813315537186321970719165779488475377688633084782731170482174374"+
			"7256947679650787426167575073363006751150073195493002048627629373"+
			"76227751462258339344895829332", 10)
	lowlvlTag := objects.NewTag(
		objects.HID(b.Hash()),
		"blob",
		"testBlob",
		nil,
		taghkid,
	) //gen test tag
	services.PostTag(lowlvlTag)
	testtag, err := services.GetTag(taghkid, lowlvlTag.NameSegment)
	if err != nil {
		t.Fatalf("GetTag Fail. error: %s", err)
	} else if !bytes.Equal(testtag.Hkid, taghkid) {
		t.Fatalf("GetTag Fail.\n\t expected: %v\n\t got: %v", taghkid, testtag.Hkid)
	} else if !testtag.Verify() {
		t.Fatalf("GetTag Verify Fail")
	}

	prikey, err := services.GetKey(taghkid)
	if err != nil || !bytes.Equal(prikey.Hkid(), taghkid) || !prikey.Verify() {
		t.Logf("GetKey Fail")
		t.Fail()
	}
}

func TestPostBlob(t *testing.T) {
	services.PostBlob(objects.Blob{})
	testhkid := objects.HkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
	testpath := "TestPostBlob"
	indata := []byte("TestPostData")

	services.Post(testhkid, testpath, objects.Blob(indata))
	outdata, err := services.Get(testhkid, testpath)
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
	testhkid := objects.HkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
	testpath := "TestPostList1/TestPostList2/TestPostBlob"
	indata := []byte("TestPostListListBlobData")
	services.Post(testhkid, testpath, objects.Blob(indata))
	outdata, err := services.Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostCommitBlob(t *testing.T) {
	testhkid := objects.HkidFromDString("65232373562705602286177837897283"+
		"2941659551264911224937349783059207224141689361121606942380473043"+
		"7860475300564272976762085068519188612732562106886379081213385", 10)
	testRepoHkid := objects.HkidFromDString("22371143209450593169269383669277"+
		"1410459232098247632372342448006863927240156318431751613873181811"+
		"3842641285036340692759591635625837820111726090732747634977413", 10)
	services.InsertRepo(testhkid, "TestPostCommit", testRepoHkid)
	testpath := "TestPostCommit/TestPostBlob"
	indata := []byte("TestPostCommitBlobData")
	services.Post(testhkid, testpath, objects.Blob(indata))
	outdata, err := services.Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostTagBlob(t *testing.T) {
	testhkid := objects.HkidFromDString("65232373562705602286177837897283"+
		"2941659551264911224937349783059207224141689361121606942380473043"+
		"7860475300564272976762085068519188612732562106886379081213385", 10)
	testDomainHkid := objects.HkidFromDString("48133155371863219707191657794884"+
		"7537768863308478273117048217437472569476796507874261675750733630"+
		"0675115007319549300204862762937376227751462258339344895829332", 10)
	services.InsertDomain(testhkid, "TestPostTag", testDomainHkid)
	testpath := "TestPostTag/TestPostBlob"
	indata := []byte("TestPostTagBlobData")
	services.Post(testhkid, testpath, objects.Blob(indata))
	outdata, err := services.Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostTagCommitBlob(t *testing.T) {
	testhkid := objects.HkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domainHkid := objects.HkidFromDString("39968110670682397993178679825250"+
		"9423226869972672234437068973021071131498376777586055610149840018"+
		"5744208447673206609026128894016152514163591905578729891874833", 10)
	repoHkid := objects.HkidFromDString("94522678075182002377842140271746"+
		"1502019766737301062946423046280258817349516439546479625226895211"+
		"80808353215150536034481206091147220911087792299373183736254", 10)
	err := services.InsertDomain(testhkid, "testTag", domainHkid)
	if err != nil {
		t.Errorf("InsertDomain with error: %s", err)
	}
	err = services.InsertRepo(testhkid, "testTag/testCommit", repoHkid)
	indata := objects.Blob([]byte("TestTagCommitBlobData"))
	testpath := "testTag/testCommit/testBlob"
	_, err = services.Post(testhkid, testpath, indata)
	outdata, err := services.Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostTagTagBlob(t *testing.T) {
	testhkid := objects.HkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domain1Hkid := objects.HkidFromDString("32076859881811206392323279987831"+
		"3732334949433938278619036381396945204532895319697233476900342324"+
		"1692155627674142765782672165943417419038514237233188152538761", 10)
	domain2Hkid := objects.HkidFromDString("49220288257701056900010210884990"+
		"0714973444364727181180850528073586453638681999434006549298762097"+
		"4197649255374796716934112121800838847071661501215957753532505", 10)
	err := services.InsertDomain(testhkid, "testTag1", domain1Hkid)
	if err != nil {
		t.Errorf("InsertDomain with error: %s", err)
	}
	err = services.InsertDomain(testhkid, "testTag1/testTag2", domain2Hkid)
	if err != nil {
		t.Errorf("InsertDomain with error: %s", err)
	}
	indata := objects.Blob([]byte("TestTagTagBlobData"))
	testpath := "testTag1/testTag2/testBlob"
	_, err = services.Post(testhkid, testpath, indata)
	outdata, err := services.Get(testhkid, testpath)
	if err != nil {
		log.Printf("[TestPostTagTagBlob] testhkid %s", testhkid)
		log.Printf("[TestPostTagTagBlob] domain1Hkid %s", domain1Hkid)
		log.Printf("[TestPostTagTagBlob] domain2Hkid %s", domain2Hkid)
		t.Errorf("Retreved with error: %s", err)
	} else if !bytes.Equal(indata, outdata) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", indata, outdata)
	}
}

func TestPostListCommitBlob(t *testing.T) {
	testhkid := objects.HkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	repoHkid := objects.HkidFromDString("59803288656043806807393139191118"+
		"0777289273091938777159029927847596771500408478956278378281366717"+
		"7487960901581583946753338859223459810645621124266443931192097", 10)
	err := services.InsertDomain(testhkid, "testList/testCommit", repoHkid)
	if err != nil {
		t.Errorf("InsertDomain with error: %s", err)
	}
	indata := objects.Blob([]byte("TestListCommitBlobData"))
	testpath := "testList/testCommit/testBlob"
	_, err = services.Post(testhkid, testpath, indata)
	outdata, err := services.Get(testhkid, testpath)
	if err != nil {
		log.Printf("[TestPostListCommitBlob] testhkid %s", testhkid)
		log.Printf("[TestPostListCommitBlob] repoHkid %s", repoHkid)
		t.Errorf("Retreved with error: %s", err)
	} else if !bytes.Equal(indata, outdata) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", indata, outdata)
	}
}

func TestPostTagListBlob(t *testing.T) {
	testhkid := objects.HkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domainHkid := objects.HkidFromDString("62089221762245310704629142682144"+
		"1944826557905230450143203631438168806532495876980559885034903315"+
		"4294997505754401230560960060918213268981906409591978967796584", 10)
	err := services.InsertDomain(testhkid, "testTag", domainHkid)
	if err != nil {
		t.Errorf("InsertDomain with error: %s", err)
	}
	indata := objects.Blob([]byte("TestTagListBlobData"))
	testpath := "testTag/testList/testBlob"
	_, err = services.Post(testhkid, testpath, indata)
	outdata, err := services.Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		t.Fail()
	}
}

func TestPostListTagBlob(t *testing.T) {
	testhkid := objects.HkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)
	domainHkid := objects.HkidFromDString("12796633883654089746486670711069"+
		"9781359720828332046318301886846633714179790444071153863702142701"+
		"146332294245448463914286494124849121460550667767568731696934", 10)
	err := services.InsertDomain(testhkid, "testList/testTag", domainHkid)
	indata := objects.Blob([]byte("TestListTagBlobData"))
	testpath := "testList/testTag/testBlob"
	_, err = services.Post(testhkid, testpath, indata)
	outdata, err := services.Get(testhkid, testpath)
	if err != nil {
		t.Errorf("Retreved with error: %s", err)
	} else if !bytes.Equal(indata, outdata) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", indata, outdata)
	}
}

func TestGetBlob(t *testing.T) {
	//hkidFromDString("4629814823893296480016411334808793836186124559723200"+
	//	"9799621767537249764640887064630013835561124248209118706504211519889067517"+
	//	"10824965155500230480521264034469", 10) //store the key
	setupForGets()

	dabytes, err := hex.DecodeString("1312ac161875b270da2ae4e1471ba94a" +
		"9883419250caa4c2f1fd80a91b37907e")

	hkid := objects.HKID(dabytes)
	path := "testTag/testBlob"
	b := []byte(":(")
	if err == nil {
		//log.Println(hkid.Hex())
		b, err = services.Get(hkid, path)
	}
	if !bytes.Equal([]byte("testing"), b) || err != nil {
		t.Fail()
	}
}

func TestGetList(t *testing.T) {
	testhkid := objects.HkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
	testpath := "TestPostList1/TestPostList2/TestPostBlob"
	indata := []byte("TestPostListListBlobData")
	services.Post(testhkid, testpath, objects.Blob(indata))
	outdata, err := services.Get(testhkid, testpath)
	if !bytes.Equal(indata, outdata) || err != nil {
		log.Printf("\nTestGetList:\nExpected: %s\nGot: %s\n", indata, outdata)
		t.Fail()
	}
}

func TestGetCommit(t *testing.T) {
	testhkid := objects.HkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
	outdata, err := services.Get(testhkid, "TestPostCommit")
	outlist, err := objects.ListFromBytes(outdata)
	typeString, hash := outlist.HashForNamesegment("TestPostBlob")
	//truthdata := []byte("90014ae279fa5034a51def77132457cd" +
	//	"66403facc3d88b54bd3e84ecade8f633,blob,TestPostBlob")
	if typeString != "blob" ||
		hash.Hex() != "90014ae279fa5034a51def77132457cd"+
			"66403facc3d88b54bd3e84ecade8f633" ||
		err != nil {
		log.Printf("\n\tTestGetList:\n\t%s\n", outdata)
		t.Fail()
	}
}

func TestGetTag(t *testing.T) {
	//think about what it means to get a domain with no path
	testhkid := objects.HkidFromDString("65232373562705602286177837897283294165955126"+
		"49112249373497830592072241416893611216069423804730437860475300564272"+
		"976762085068519188612732562106886379081213385", 10)
	outdata, err := services.Get(testhkid, "TestPostTag")
	truthdata := []byte("")
	if !bytes.Equal(truthdata, outdata) {
		log.Printf("\n\tTestGetList:\n\t%s\n\terror: %s\n", outdata, err)
		t.Fail()
	}
}

func TestKeyGen(t *testing.T) {
	t.SkipNow()
	priv := objects.KeyGen()
	log.Printf("TestKeyGen\nX = %v\nY = %v\nD = %v\n", priv.PublicKey.X,
		priv.PublicKey.Y, priv.D)
	err := services.PostKey(priv)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	services.PostBlob(elliptic.Marshal(priv.PublicKey.Curve,
		priv.PublicKey.X, priv.PublicKey.Y))
}

func Testbadfrombytes(t *testing.T) {
	var err error
	_, err = objects.ListFromBytes([]byte{})
	if err.Error() != "Could not parse list bytes" {
		t.Errorf("[] should not parse")
	}
	_, err = objects.CommitFromBytes([]byte{})
	if err.Error() != "Could not parse commit bytes" {
		t.Errorf("[] should not parse")
	}
	_, err = objects.TagFromBytes([]byte{})
	if err.Error() != "Could not parse tag bytes" {
		t.Errorf("[] should not parse")
	}
	_, err = objects.PrivateKeyFromBytes([]byte{})
	if err.Error() != "Could not parse commit bytes" {
		t.Errorf("[] should not parse")
	}
}

func setupForGets() {
	hkidT := objects.HkidFromDString("39968110670682397993178679825250942322686997267223"+
		"4437068973021071131498376777586055610149840018574420844767320660902612889"+
		"4016152514163591905578729891874833", 10)

	//key for commit
	hkidC := objects.HkidFromDString("4629814823893296480016411334808793836186124559723200"+
		"9799621767537249764640887064630013835561124248209118706504211519889067517"+
		"10824965155500230480521264034469", 10)

	//Post blob
	testBlob := objects.Blob([]byte("testing")) //gen test blob
	err := services.PostBlob(testBlob)          //store test blob
	if err != nil {
		log.Println(err)
	}

	//post tag
	testTagPointingToTestBlob := objects.NewTag(
		objects.HID(testBlob.Hash()),
		"blob",
		"testBlob",
		nil,
		hkidT,
	) //gen test tag
	err = services.PostTag(testTagPointingToTestBlob) //post test tag
	if err != nil {
		log.Println(err)
	}

	//post list
	testListPiontingToTestTag := objects.NewList(testTagPointingToTestBlob.Hkid,
		"tag",
		"testTag") //gen test list
	err = services.PostBlob(testListPiontingToTestTag.Bytes()) //store test list
	if err != nil {
		log.Println(err)
	}

	// post commit
	testCommitPointingToTestList := objects.NewCommit(testListPiontingToTestTag.Hash(),
		hkidC) //gen test commit
	err = services.PostCommit(testCommitPointingToTestList) //post test commit
	if err != nil {
		log.Println(err)
	}
}
