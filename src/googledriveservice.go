package main

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/drive/v2"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var blobsFolderId string
var commitsFolderId string
var tagsFolderId string
var keysFolderId string
var transport *oauth.Transport
var driveService *drive.Service

func init() { // Set up a configuration.
	log.SetFlags(log.Lshortfile)

	config := &oauth.Config{
		ClientId: "755660992417.apps.googleusercontent.com",
		// from https://code.google.com/apis/console/
		//ClientSecret: "h8bA_4cKRD8nSE6kzC9vAEw2",
		ClientSecret: "",
		//TODO move this out of code and Reset client secret...
		// from https://code.google.com/apis/console/
		RedirectURL: "oob",
		Scope:       drive.DriveReadonlyScope,
		AuthURL:     "https://accounts.google.com/o/oauth2/auth",
		TokenURL:    "https://accounts.google.com/o/oauth2/token",
		TokenCache:  oauth.CacheFile("../bin/tokencachefile.json"),
	}

	//code := "4/8JG6IFW1v7QucMzD-aaTG4TmqpjO.khB140k6OnMfshQV0ieZDAqGXAxGgwI"
	code := ""

	// Set up a Transport using the config.
	transport = &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			log.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			log.Println(url)
			panic("")
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		log.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token

	// Make the request.
	httpClient := transport.Client()
	r, err := httpClient.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer r.Body.Close()

	// Write the response to standard output.
	//log.Println(r.Body)
	driveService, err = drive.New(httpClient)
	ccfsFolderId, err := getChildWithTitle(driveService, "root", "ccfs")
	blobsFolderId, err = getChildWithTitle(driveService, ccfsFolderId, "blobs")
	commitsFolderId, err = getChildWithTitle(driveService, ccfsFolderId, "commits")
	tagsFolderId, err = getChildWithTitle(driveService, ccfsFolderId, "tags")
	keysFolderId, err = getChildWithTitle(driveService, ccfsFolderId, "keys")
	log.Printf("\n\tccfsFolderId: %v\n\tblobsFolderId: %v\n\t"+
		"commitsFolderId: %v\n\ttagsFolderId: %v\n\tkeysFolderId: %v\n",
		ccfsFolderId, blobsFolderId, commitsFolderId,
		tagsFolderId, keysFolderId)

	hb, err := HcidFromHex("42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12")
	b, err := googledriveservice_GetBlob(hb)
	log.Printf("\n\tBlob Contents: %v \n\tError: %v",
		b.Hash().String() == hb.Hex(), err)

	hc, err := HkidFromHex("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02")
	c, err := googledriveservice_GetCommit(hc)
	log.Printf("\n\tCommit Contents: %v \n\tError: %v", c.Verifiy(), err)

	ht, err := HkidFromHex("f65b92b9ce15e167b98fc896f0a365c87c39565642a59ba0060db3b33be6d885")
	t, err := googledriveservice_GetTag(ht, "testBlob")
	log.Printf("\n\tTag Contents: %v \n\tError: %v", t.Verifiy(), err)

	//fileId, err := getChildWithTitle(driveService, blobsFolderId,
	//	"42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12")
	//f, err := driveService.Files.Get(fileId).Do()
	//fileString, err := DownloadFile(driveService, transport, f)
	//if err != nil {
	//	log.Printf("An error occurred: %v\n", err)
	//	return
	//} else {
	//	log.Printf("File Contents: %s\n", fileString)
	//}

	//log.Printf("Title: %v", f.Title)
	//log.Printf("Description: %v", f.Description)
	//log.Printf("MIME type: %v", f.MimeType)
	//log.Printf("Download Url: %v", f.DownloadUrl)

	//var cs []*drive.ChildReference
	//cs, err = AllChildren(driveService, fileId,
	//  	"title = '42cc3a4c4a9d9d3ee7de9322b45acb0e5a5c33550d9ad4791df6ae937a869e12'")
	//log.Printf("len(Children): %v", len(cs))
	//for _, v := range cs {
	//	f, err := driveService.Files.Get(v.Id).Do()
	//	if err != nil {
	//		log.Printf("Children: %v", f.Title)
	//	} else {
	//		log.Printf("Get Error: %v", f.Id)
	//	}
	//}

	log.Printf("Done:")
}

func getChildWithTitle(d *drive.Service, parentId string, title string) (string, error) {
	r, err := d.Children.List(parentId).Q(fmt.Sprintf("title = '%s'", title)).Do()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if len(r.Items) < 1 {
		return "", fmt.Errorf("no file with title %s", title)
	}
	return r.Items[0].Id, err
}

// AllChildren fetches all the children of a given folder
func AllChildren(d *drive.Service, folderId string, qString string) ([]*drive.ChildReference,
	error) {
	var cs []*drive.ChildReference
	pageToken := ""
	for {
		q := d.Children.List(folderId).Q(qString)
		// If we have a pageToken set, apply it to the query
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}
		r, err := q.Do()
		if err != nil {
			log.Printf("An error occurred: %v\n", err)
			return cs, err
		}
		cs = append(cs, r.Items...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return cs, nil
}

// DownloadFile downloads the content of a given file object
func DownloadFile(d *drive.Service, t http.RoundTripper, f *drive.File) ([]byte, error) {
	// t parameter should use an oauth.Transport
	downloadUrl := f.DownloadUrl
	if downloadUrl == "" {
		// If there is no downloadUrl, there is no body
		fmt.Printf("An error occurred: File is not downloadable")
		return []byte{}, nil
	}
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return []byte{}, err
	}
	resp, err := t.RoundTrip(req)
	// Make sure we close the Body later
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return []byte{}, err
	}
	return body, nil
}

func googledriveservice_GetBlob(hash HCID) (b blob, err error) {
	fileId, err := getChildWithTitle(driveService, blobsFolderId,
		hash.Hex())
	f, err := driveService.Files.Get(fileId).Do()
	fileString, err := DownloadFile(driveService, transport, f)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	return blob(fileString), err
}

func googledriveservice_GetCommit(hash HKID) (c commit, err error) {
	thisCommitFolderId, err := getChildWithTitle(driveService, commitsFolderId,
		hash.Hex())
	r, err := driveService.Children.List(thisCommitFolderId).Do()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if len(r.Items) < 1 {
		return commit{}, fmt.Errorf("no file %s", hash.Hex())
	}
	latestTitle := ""
	thisCommitfile := new(drive.File)
	for _, item := range r.Items {
		f, err := driveService.Files.Get(item.Id).Do()
		if f.Title > latestTitle && err == nil {
			latestTitle = f.Title
			thisCommitfile = f
		}

	}
	//thisCommitID, err := getChildWithTitle(driveService,
	//	thisCommitFolderId, latestID)
	//thisCommitfile, err := driveService.Files.Get(thisCommitID).Do()
	commitBytes, err := DownloadFile(driveService, transport, thisCommitfile)
	c, err = CommitFromBytes(commitBytes)
	return c, err
}

func googledriveservice_GetKey(hash HKID) (data blob, err error) {
	fileId, err := getChildWithTitle(driveService, keysFolderId,
		hash.Hex())
	f, err := driveService.Files.Get(fileId).Do()
	fileString, err := DownloadFile(driveService, transport, f)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	return blob(fileString), err
}

func googledriveservice_GetTag(hash HKID, nameSegment string) (t Tag, err error) {
	hkidTagFolderId, err := getChildWithTitle(driveService, tagsFolderId,
		hash.Hex())
	nameSegmentTagFolderId, err := getChildWithTitle(driveService, hkidTagFolderId,
		nameSegment)
	r, err := driveService.Children.List(nameSegmentTagFolderId).Do()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if len(r.Items) < 1 {
		return Tag{}, fmt.Errorf("no file %s ,%s", hash.Hex(), nameSegment)
	}
	latestTitle := ""
	thisTagfile := new(drive.File)
	for _, item := range r.Items {
		f, err := driveService.Files.Get(item.Id).Do()
		if f.Title > latestTitle && err == nil {
			latestTitle = f.Title
			thisTagfile = f
		}

	}
	tagBytes, err := DownloadFile(driveService, transport, thisTagfile)
	t, err = TagFromBytes(tagBytes)
	return t, nil
}
