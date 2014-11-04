// benchmark_test.go
package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/localfile"
)

var benchmarkRepo objects.HKID
var benchmarkCommitHkid objects.HKID
var benchmarkTagHkid objects.HKID

func init() {

	//Open|Create Log File
	logFileName := "bin/log.txt"
	logFile, err := os.OpenFile(logFileName, os.O_RDWR, 664)
	if err != nil {
		log.Println("Unable to Create/Open log file. No output will be captured.")
	} else {
		log.SetOutput(logFile)
	}
	log.SetFlags(log.Lshortfile)

	objects.RegisterGeterPoster(
		services.GetPublicKeyForHkid,
		services.GetPrivateKeyForHkid,
		services.PostKey,
		services.PostBlob,
	)
	localfile.Start()
	//Post(benchmarkRepo, "BlobFound", blob("blob found"))
	//Post(benchmarkRepo, "listFound/BlobFound", blob("list found"))
	benchmarkRepo = objects.HkidFromDString("44089867384569081480066871308647"+
		"4832666868594293316444099156169623352946493325312681245061254048"+
		"6538169821270508889792789331438131875225590398664679212538621", 10)
	log.Printf("Benchmark Repository HKID: %s\n", benchmarkRepo)
	benchmarkCommitHkid = objects.HkidFromDString("36288652923287342545336063204999"+
		"9357791979761632757400493952327434464825857894440491353330036559"+
		"1539025688752776406270441884985963379175226110071953813093104", 10)
	log.Printf("Benchmark Commit HKID: %s\n", benchmarkCommitHkid)
	benchmarkTagHkid = objects.HkidFromDString("54430439211086161065670118078952"+
		"6855263811485554809970416168964497131310494084201908881724207840"+
		"2305843436117034888111308798569392135240661266075941854101839", 10)
	log.Printf("Benchmark Tag HKID: %s\n", benchmarkTagHkid)
}

//hkid 549baa6497db3615332aae859680b511117e299879ee311fbac4d1a40f93b8d0

func BenchmarkStoreOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := os.MkdirAll("bin/", 0764)
		err = ioutil.WriteFile("bin/storeone", []byte("storeone"), 0664)
		if err != nil {
			log.Panic(err)
		}
		_, err = ioutil.ReadFile("bin/storeone")
		if err != nil {
			log.Panic(err)
		}
		//log.Print(data)
	}
}

func BenchmarkLowLevelPath(b *testing.B) {
	//key for tag
	hkidT := objects.HkidFromDString("39968110670682397993178679825250942322686997267223"+
		"4437068973021071131498376777586055610149840018574420844767320660902612889"+
		"4016152514163591905578729891874833", 10)

	//key for commit
	hkidC := objects.HkidFromDString("4629814823893296480016411334808793836186124559723200"+
		"9799621767537249764640887064630013835561124248209118706504211519889067517"+
		"10824965155500230480521264034469", 10)

	for i := 0; i < b.N; i++ {

		blobHcid := postBlob("testing")
		postTag(blobHcid, hkidT)
		tagHcid := postList(hkidT)
		postCommit(tagHcid, hkidC)

		//get commit
		hkid, _ := hex.DecodeString(
			"1312ac161875b270da2ae4e1471ba94a9883419250caa4c2f1fd80a91b37907e",
		)
		testcommit, err := services.GetCommit(hkid)
		if err != nil {
			log.Panic(err)
		}

		//log.Printf("authentic commit:%v\n", testcommit.Verifiy())
		if !testcommit.Verify() {
			b.FailNow()
		}
		//get list
		listbytes, err := services.GetBlob(testcommit.ListHash)
		if err != nil {
			log.Panic(err)
		}

		testlist, err := objects.ListFromBytes(listbytes)

		//log.Printf("authentic list:%v\n", )
		if !bytes.Equal(testcommit.ListHash, testlist.Hash()) {
			b.FailNow()
		}
		//get tag
		_, testTagHash := testlist.HashForNamesegment("testTag")
		testTag, err := services.GetTag(testTagHash.Bytes(), "testBlob")
		//log.Printf("authentic tag:%v\n", testTag.Verifiy())
		if !testTag.Verify() {
			b.FailNow()
		}
		//get blob
		testBlob, err := services.GetBlob(testTag.HashBytes.(objects.HCID))
		//log.Printf("authentic blob:%v\n", bytes.Equal(testTag.HashBytes,
		//	testBlob.Hash()))
		if err != nil {
			b.Fatalf("Get Blob Error: %s\n", err)
		}
		if !bytes.Equal(testTag.HashBytes.(objects.HCID), testBlob.Hash()) {
			b.FailNow()
		}
	}
}

func BenchmarkHighLevelPath(b *testing.B) {
	testhkid := objects.HkidFromDString("46298148238932964800164113348087"+
		"9383618612455972320097996217675372497646408870646300138355611242"+
		"4820911870650421151988906751710824965155500230480521264034469", 10)

	indata := objects.Blob([]byte("testing"))
	testpath := "testTag/testBlob"
	//_, err = Post(testhkid, testpath, indata)
	for i := 0; i < b.N; i++ {
		outdata, err := services.Get(testhkid, testpath)
		if err != nil {
			b.Fatalf("Failed to retrieve Blob: %s", err)
		}
		if !bytes.Equal(indata, outdata) {
			b.Fatalf("\tExpected: %s\n\tActual: %s\n\tErr:%s\n",
				indata, outdata, err)
		}
	}
}

//BenchmarkBlobFound times the retrieval of a blob that is found
func BenchmarkBlobFound(b *testing.B) {
	indata := objects.Blob("Blob Found Data")
	services.Post(benchmarkRepo, "BlobFound", indata)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "BlobFound")
		if err != nil {
			b.Fatalf("Failed to find Blob: %s", err)
		}
	}
}

//BenchmarkBlobNotFound times the retrieval of a blob that is not found
func BenchmarkBlobNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "BlobNotFound")
		if err.Error() != "Blob not found" {
			b.Fatalf("Blob incorrectly detected: %s", err)
		}
	}
}

//BenchmarkBlobInsert times the posting of a blob to a repository.
func BenchmarkBlobInsert(b *testing.B) {
	indata := objects.Blob("Benchmark Blob Data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services.Post(benchmarkRepo, "benchmarkBlob", indata)
		//TODO clear blob after each insert
	}
	outdata, err := services.Get(benchmarkRepo, "benchmarkBlob")
	if !bytes.Equal(indata, outdata) || err != nil {
		b.Fatalf("\tExpected: %s\n\tActual: %s\n\tErr:%s\n",
			indata, outdata, err)
	}
}

func BenchmarkBlobUpdate(b *testing.B) {
	indata := objects.Blob("Benchmark Blob Data")
	services.Post(benchmarkRepo, "benchmarkBlob", indata) //  to ensure is update
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services.Post(benchmarkRepo, "benchmarkBlob", indata)
	}
	outdata, err := services.Get(benchmarkRepo, "benchmarkBlob")
	if !bytes.Equal(indata, outdata) || err != nil {
		b.Fatalf("\tExpected: %s\n\tActual: %s\n\tErr:%s\n",
			indata, outdata, err)
	}
}

func BenchmarkListBlobFound(b *testing.B) {
	indata := objects.Blob("List Blob Found Data")
	services.Post(benchmarkRepo, "listFound/BlobFound", indata)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "listFound/BlobFound")
		if err != nil {
			b.Fatalf("Failed to find List Blob: %s", err)
		}
	}
}

func BenchmarkListBlobNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "listNotFound/BlobFound")
		if err.Error() != "Blob not found" {
			b.Fatalf("List Blob incorrectly detected: %s", err)
		}
	}
}

func BenchmarkListBlobInsert(b *testing.B) {
	indata := objects.Blob("Benchmark Blob Data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services.Post(benchmarkRepo, "listFound/benchmarkBlob", indata)
		//TODO clear list after each insert
	}
	outdata, err := services.Get(benchmarkRepo, "listFound/benchmarkBlob")
	if !bytes.Equal(indata, outdata) || err != nil {
		b.Fatalf("\tExpected: %s\n\tActual: %s\n\tErr:%s\n",
			indata, outdata, err)
	}
}

func BenchmarkListBlobUpdate(b *testing.B) {
	indata := objects.Blob("Benchmark Blob Data")
	services.Post(benchmarkRepo, "listFound/benchmarkBlob", indata)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services.Post(benchmarkRepo, "listFound/benchmarkBlob", indata)
	}
	outdata, err := services.Get(benchmarkRepo, "listFound/benchmarkBlob")
	if !bytes.Equal(indata, outdata) || err != nil {
		b.Fatalf("\tExpected: %s\n\tActual: %s\n\tErr:%s\n",
			indata, outdata, err)
	}
}

func BenchmarkCommitBlobFound(b *testing.B) {
	err := services.InsertRepo(benchmarkRepo, "commitFound", benchmarkCommitHkid)
	if err != nil {
		b.Fatalf("Unable to insert Commit: %s", err)
	}
	_, err = services.Post(benchmarkRepo, "commitFound/benchmarkBlob",
		objects.Blob("Benchmark Blob Data"))
	if err != nil {
		b.Fatalf("Unable to post Commit Blob: %s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "commitFound/benchmarkBlob")
		if err != nil {
			b.Fatalf("Unable to retrieve Commit Blob: %s", err)
		}
	}
}

func BenchmarkCommitBlobNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "commitNotFound/benchmarkBlob")
		if err.Error() != "Blob not found" {
			b.Fatal(err)
		}
	}
}

func BenchmarkCommitBlobInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := services.InsertRepo(benchmarkRepo, "commitFound", benchmarkCommitHkid)
		//TODO clear list after each insert
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCommitBlobUpdate(b *testing.B) {
	err := services.InsertRepo(benchmarkRepo, "commitFound", benchmarkCommitHkid)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = services.InsertRepo(benchmarkRepo, "commitFound", benchmarkCommitHkid)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTagBlobFound(b *testing.B) {
	err := services.InsertDomain(benchmarkRepo, "tagFound", benchmarkTagHkid)
	if err != nil {
		b.Fatal(err)
	}
	_, err = services.Post(benchmarkRepo, "tagFound/benchmarkBlob",
		objects.Blob("Benchmark Blob Data"))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "tagFound/benchmarkBlob")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTagBlobNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := services.Get(benchmarkRepo, "tagNotFound/benchmarkBlob")
		if err.Error() != "Blob not found" {
			b.Fatal(err)
		}
	}
}

func BenchmarkTagBlobInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := services.InsertDomain(benchmarkRepo, "tagFound", benchmarkTagHkid)
		//TODO clear list after each insert
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTagBlobUpdate(b *testing.B) {
	err := services.InsertRepo(benchmarkRepo, "tagFound", benchmarkTagHkid)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = services.InsertDomain(benchmarkRepo, "tagFound", benchmarkTagHkid)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func postBlob(data string) objects.HCID {
	testBlob := objects.Blob([]byte(data)) //gen test blob
	err := services.PostBlob(testBlob)     //store test blob
	if err != nil {
		log.Println(err)
	}
	return testBlob.Hash()

}

func postTag(objHash objects.HCID, tagHkid objects.HKID) {
	testTagPointingToTestBlob := objects.NewTag(
		objects.HID(objHash),
		"blob",
		"testBlob",
		nil,
		tagHkid,
	) //gen test tag
	err := services.PostTag(testTagPointingToTestBlob) //post test tag
	if err != nil {
		log.Println(err)
	}
}

func postList(target objects.HKID) (testTagHash objects.HCID) {
	testListPointingToTestTag := objects.NewList(target,
		"tag",
		"testTag") //gen test list
	err := services.PostList(testListPointingToTestTag) //store test list
	if err != nil {
		log.Println(err)
	}
	return testListPointingToTestTag.Hash()
}

func postCommit(hcid objects.HCID, hkidC objects.HKID) {
	testCommitPointingToTestList := objects.NewCommit(hcid,
		hkidC) //gen test commit
	err := services.PostCommit(testCommitPointingToTestList) //post test commit
	if err != nil {
		log.Println(err)
	}
}
