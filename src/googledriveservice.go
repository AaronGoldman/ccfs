package main

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/drive/v2"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (gds googledriveservice) getChildWithTitle(parentId string, title string) (string, error) {
	if gds.driveService == nil {
		log.Println("drive.Service not initialized")
		return "", fmt.Errorf("drive.Service not initialized")
	}
	//log.Println(parentId, "\n\t", title)
	r, err := gds.driveService.Children.List(parentId).Q(fmt.Sprintf("title = '%s'", title)).Do()
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
func (gds googledriveservice) AllChildren(d *drive.Service, folderId string, qString string) ([]*drive.ChildReference,
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
func (gds googledriveservice) DownloadFile(f *drive.File) ([]byte, error) {
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
	blobsFolderId   string
	commitsFolderId string
	tagsFolderId    string
	keysFolderId    string
	transport       *oauth.Transport
	driveService    *drive.Service
}

func (gds googledriveservice) getBlob(h HCID) (blob, error) {
	if gds.driveService == nil {
		return nil, fmt.Errorf("Drive Service not initialized")
	}
	fileId, err := gds.getChildWithTitle(gds.blobsFolderId, h.Hex())
	f, err := gds.driveService.Files.Get(fileId).Do()
	fileString, err := gds.DownloadFile(f)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	return blob(fileString), err
}
func (gds googledriveservice) getCommit(h HKID) (c commit, err error) {
	if gds.driveService == nil {
		return commit{}, fmt.Errorf("Drive Service not initialized")
	}
	thisCommitFolderId, err := gds.getChildWithTitle(gds.commitsFolderId, h.Hex())
	r, err := gds.driveService.Children.List(thisCommitFolderId).Do()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if len(r.Items) < 1 {
		return commit{}, fmt.Errorf("no file %s", h.Hex())
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
	c, err = CommitFromBytes(commitBytes)
	return c, err
}
func (gds googledriveservice) getTag(h HKID, namesegment string) (t tag, err error) {
	if gds.driveService == nil {
		return tag{}, fmt.Errorf("Drive Service not initialized")
	}
	hkidTagFolderId, err := gds.getChildWithTitle(gds.tagsFolderId, h.Hex())
	nameSegmentTagFolderId, err := gds.getChildWithTitle(hkidTagFolderId, namesegment)
	r, err := gds.driveService.Children.List(nameSegmentTagFolderId).Do()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return tag{}, err
	}
	if len(r.Items) < 1 {
		return tag{}, fmt.Errorf("no file %s ,%s", h.Hex(), namesegment)
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
	t, err = TagFromBytes(tagBytes)
	return t, nil
}
func (gds googledriveservice) getKey(h HKID) (b blob, err error) {
	if gds.driveService == nil {
		return nil, fmt.Errorf("Drive Service not initialized")
	}
	fileId, err := gds.getChildWithTitle(gds.keysFolderId, h.Hex())
	f, err := gds.driveService.Files.Get(fileId).Do()
	fileString, err := gds.DownloadFile(f)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	return blob(fileString), err
}

func googledriveserviceFactory() googledriveservice {
	log.SetFlags(log.Lshortfile)
	gds := googledriveservice{}
	// Set up a configuration.
	config := &oauth.Config{
		ClientId: "755660992417.apps.googleusercontent.com",
		// from https://code.google.com/apis/console/
		ClientSecret: "h8bA_4cKRD8nSE6kzC9vAEw2",
		//TODO move this out of code and Reset client secret...
		// from https://code.google.com/apis/console/
		RedirectURL: "oob",
		Scope:       drive.DriveReadonlyScope,
		AuthURL:     "https://accounts.google.com/o/oauth2/auth",
		TokenURL:    "https://accounts.google.com/o/oauth2/token",
		TokenCache:  oauth.CacheFile("../bin/tokencachefile.json"),
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
	ccfsFolderId, err := gds.getChildWithTitle("root", "ccfs")
	gds.blobsFolderId, err = gds.getChildWithTitle(ccfsFolderId, "blobs")
	gds.commitsFolderId, err = gds.getChildWithTitle(ccfsFolderId, "commits")
	gds.tagsFolderId, err = gds.getChildWithTitle(ccfsFolderId, "tags")
	gds.keysFolderId, err = gds.getChildWithTitle(ccfsFolderId, "keys")
	return gds
}

var googledriveserviceInstance googledriveservice

func dont_init() {
	googledriveserviceInstance = googledriveserviceFactory()
	log.Printf(
		"\n\tblobsFolderId: %v"+
			"\n\tcommitsFolderId: %v"+
			"\n\ttagsFolderId: %v"+
			"\n\tkeysFolderId: %v"+
			"\n\tdriveService: %v"+
			"\n\ttransport: %v\n",
		googledriveserviceInstance.blobsFolderId,
		googledriveserviceInstance.commitsFolderId,
		googledriveserviceInstance.tagsFolderId,
		googledriveserviceInstance.keysFolderId,
		googledriveserviceInstance.driveService,
		googledriveserviceInstance.transport,
	)
}
