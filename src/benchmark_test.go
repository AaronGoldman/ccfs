// benchmark_test.go
package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var setup_for_benchmark func() (HKID, HKID, HKID)

func init() {
	setup_for_benchmark = func() (HKID, HKID, HKID) {
		log.SetFlags(log.Lshortfile)
		benchmarkRepo := hkidFromDString("44089867384569081480066871308647"+
			"4832666868594293316444099156169623352946493325312681245061254048"+
			"6538169821270508889792789331438131875225590398664679212538621", 10)
		benchmarkCommitHkid := hkidFromDString(
			"362886529232873425453360632049999357791979761632757400493952327"+
				"434464825857894440491353330036559153902568875277640627044188498"+
				"5963379175226110071953813093104", 10)
		benchmarkTagHkid := hkidFromDString("54430439211086161065670118078952"+
			"6855263811485554809970416168964497131310494084201908881724207840"+
			"2305843436117034888111308798569392135240661266075941854101839", 10)
		Post(benchmarkRepo, "BlobFound", blob("blob found"))
		Post(benchmarkRepo, "listFound/BlobFound", blob("list found"))

		setup_for_benchmark = func() (HKID, HKID, HKID) {
			return benchmarkRepo, benchmarkCommitHkid, benchmarkTagHkid
		}
		return benchmarkRepo, benchmarkCommitHkid, benchmarkTagHkid
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

func BenchmarkLowLevelPath(b *testing.B) {
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

func BenchmarkHighLevelPath(b *testing.B) {
	log.SetFlags(log.Lshortfile)
	testhkid := hkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)

	indata := blob([]byte("testing"))
	testpath := "testTag/testBlob"
	//_, err = Post(testhkid, testpath, indata)
	for i := 0; i < b.N; i++ {
		outdata, err := Get(testhkid, testpath)
		if !bytes.Equal(indata, outdata) || err != nil {
			log.Println(testhkid, outdata)
			b.FailNow()
		}
	}
}

func BenchmarkBlobFound(b *testing.B) {
	benchmarkRepo, _, _ := setup_for_benchmark()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(benchmarkRepo, "BlobFound")
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkBlobNotFound(b *testing.B) {
	benchmarkRepo, _, _ := setup_for_benchmark()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(benchmarkRepo, "BlobNotFound")
		if err.Error() != "Blob not found" {
			log.Println(err)
			b.FailNow()
		}
	}
}

func BenchmarkBlobInsert(b *testing.B) {
	setup_for_benchmark()
	log.SetFlags(log.Lshortfile)
	b.FailNow()
}

func BenchmarkBlobUpdate(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkListBlobFound(b *testing.B) {
	benchmarkRepo, _, _ := setup_for_benchmark()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(benchmarkRepo, "listFound/BlobFound")
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkListBlobNotFound(b *testing.B) {
	benchmarkRepo, _, _ := setup_for_benchmark()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(benchmarkRepo, "listNotFound/BlobFound")
		if err.Error() != "Blob not found" {
			log.Println(err)
			b.FailNow()
		}
	}
}

func BenchmarkListBlobInsert(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkListBlobUpdate(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkCommitBlobFound(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkCommitBlobNotFound(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkCommitBlobInsert(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkCommitBlobUpdate(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkTagBlobFound(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkTagBlobNotFound(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkTagBlobInsert(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}

func BenchmarkTagBlobUpdate(b *testing.B) {
	setup_for_benchmark()
	b.FailNow()
}
