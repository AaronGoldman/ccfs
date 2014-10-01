//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
package googledrive

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/drive/v2"
	"github.com/AaronGoldman/ccfs/objects"
)

func (gds googledriveservice) GetId() string {
	return "googledrive"
}

func (gds googledriveservice) getChildWithTitle(parentID string, title string) (string, error) {
	if gds.driveService == nil {
		log.Println("drive.Service not initialized")
		return "", fmt.Errorf("drive.Service not initialized")
	}
	//log.Println(parentID, "\n\t", title)
	r, err := gds.driveService.Children.List(parentID).Q(fmt.Sprintf("title = '%s'", title)).Do()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return "", err
	}
	if len(r.Items) < 1 {
		return "", fmt.Errorf("no file with title %s", title)
	}
	return r.Items[0].Id, err
}

// AllChildren fetches all the children of a given folder
func (gds googledriveservice) AllChildren(d *drive.Service, folderID string, qString string) ([]*drive.ChildReference,
	error) {
	var cs []*drive.ChildReference
	pageToken := ""
	for {
		q := d.Children.List(folderID).Q(qString)
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
func (gds googledriveservice) DownloadFile(f *drive.File) ([]byte, error) {
	// t parameter should use an oauth.Transport
	downloadURL := f.DownloadUrl
	if downloadURL == "" {
		// If there is no downloadUrl, there is no body
		fmt.Printf("An error occurred: File is not downloadable")
		return []byte{}, nil
	}
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return []byte{}, err
	}
	resp, err := gds.transport.RoundTrip(req)
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

type googledriveservice struct {
	blobsFolderID   string
	commitsFolderID string
	tagsFolderID    string
	keysFolderID    string
	transport       *oauth.Transport
	driveService    *drive.Service
}

func (gds googledriveservice) GetBlob(h objects.HCID) (objects.Blob, error) {
	if gds.driveService == nil {
		return nil, fmt.Errorf("Drive Service not initialized")
	}
	fileID, err := gds.getChildWithTitle(gds.blobsFolderID, h.Hex())
	f, err := gds.driveService.Files.Get(fileID).Do()
	fileString, err := gds.DownloadFile(f)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	return objects.Blob(fileString), err
}
func (gds googledriveservice) GetCommit(h objects.HKID) (c objects.Commit, err error) {
	if gds.driveService == nil {
		return objects.Commit{}, fmt.Errorf("Drive Service not initialized")
	}
	thisCommitFolderID, err := gds.getChildWithTitle(gds.commitsFolderID, h.Hex())
	r, err := gds.driveService.Children.List(thisCommitFolderID).Do()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if len(r.Items) < 1 {
		return objects.Commit{}, fmt.Errorf("no file %s", h.Hex())
	}
	latestTitle := ""
	thisCommitfile := new(drive.File)
	for _, item := range r.Items {
		f, err := gds.driveService.Files.Get(item.Id).Do()
		if err != nil {
			return c, err
		}
		if f.Title > latestTitle && err == nil {
			latestTitle = f.Title
			thisCommitfile = f
		}
	}
	commitBytes, err := gds.DownloadFile(thisCommitfile)
	c, err = objects.CommitFromBytes(commitBytes)
	return c, err
}
func (gds googledriveservice) GetTag(h objects.HKID, namesegment string) (t objects.Tag, err error) {
	if gds.driveService == nil {
		return objects.Tag{}, fmt.Errorf("Drive Service not initialized")
	}
	hkidTagFolderID, err := gds.getChildWithTitle(gds.tagsFolderID, h.Hex())
	nameSegmentTagFolderID, err := gds.getChildWithTitle(hkidTagFolderID, namesegment)
	r, err := gds.driveService.Children.List(nameSegmentTagFolderID).Do()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return objects.Tag{}, err
	}
	if len(r.Items) < 1 {
		return objects.Tag{}, fmt.Errorf("no file %s ,%s", h.Hex(), namesegment)
	}
	latestTitle := ""
	thisTagfile := new(drive.File)
	for _, item := range r.Items {
		f, err := gds.driveService.Files.Get(item.Id).Do()
		if err != nil {
			return t, err
		} //log.Println(f.Title)
		if f.Title > latestTitle && err == nil {
			latestTitle = f.Title
			thisTagfile = f
		}

	}
	tagBytes, err := gds.DownloadFile(thisTagfile)
	t, err = objects.TagFromBytes(tagBytes)
	return t, nil
}
func (gds googledriveservice) GetKey(h objects.HKID) (b objects.Blob, err error) {
	if gds.driveService == nil {
		return nil, fmt.Errorf("Drive Service not initialized")
	}
	fileID, err := gds.getChildWithTitle(gds.keysFolderID, h.Hex())
	f, err := gds.driveService.Files.Get(fileID).Do()
	fileString, err := gds.DownloadFile(f)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	return objects.Blob(fileString), err
}

func googledriveserviceFactory() googledriveservice {
	log.SetFlags(log.Lshortfile)
	gds := googledriveservice{}
	// Set up a configuration.
	config := &oauth.Config{
		ClientId: "755660992417.apps.googleusercontent.com",
		// from https://code.google.com/apis/console/
		ClientSecret: "h8bA_4cKRD8nSE6kzC9vAEw2",
		//TODO move this out of code and Reset client secret
		// from https://code.google.com/apis/console/
		RedirectURL: "oob",
		Scope:       drive.DriveReadonlyScope,
		AuthURL:     "https://accounts.google.com/o/oauth2/auth",
		TokenURL:    "https://accounts.google.com/o/oauth2/token",
		TokenCache:  oauth.CacheFile("bin/tokencachefile.json"),
	}

	code := "4/rSyLcOy_oBllG65sojDydzbxLp06.AgeuzdzuK-IWshQV0ieZDArWsFLjhAI"

	// Set up a Transport using the config.
	gds.transport = &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			log.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE")
			log.Println(url)
			panic("")
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = gds.transport.Exchange(code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		log.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	gds.transport.Token = token

	// Make the request.
	//httpClient := transport.Client()
	//r, err := httpClient.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	//if err != nil {
	//	log.Fatal("Get:", err)
	//}
	//defer r.Body.Close()

	//get the ID's of the object folders'
	httpClient := gds.transport.Client()
	gds.driveService, err = drive.New(httpClient)
	ccfsFolderID, err := gds.getChildWithTitle("root", "ccfs")
	gds.blobsFolderID, err = gds.getChildWithTitle(ccfsFolderID, "blobs")
	gds.commitsFolderID, err = gds.getChildWithTitle(ccfsFolderID, "commits")
	gds.tagsFolderID, err = gds.getChildWithTitle(ccfsFolderID, "tags")
	gds.keysFolderID, err = gds.getChildWithTitle(ccfsFolderID, "keys")

	log.Printf(
		"\n\tblobsFolderId: %v"+
			"\n\tcommitsFolderId: %v"+
			"\n\ttagsFolderId: %v"+
			"\n\tkeysFolderId: %v"+
			"\n\tdriveService: %v"+
			"\n\ttransport: %v\n",
		gds.blobsFolderID,
		gds.commitsFolderID,
		gds.tagsFolderID,
		gds.keysFolderID,
		gds.driveService,
		gds.transport,
	)

	return gds
}

//Instance is the instance of the googledriveservice
var Instance googledriveservice

func init() {
	//googledriveserviceInstance = googledriveserviceFactory()
}
