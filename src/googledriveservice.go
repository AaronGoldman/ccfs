package main

import (
	"log"
	"net/http"

	drive "code.google.com/p/google-api-go-client/drive/v2"
	//"code.google.com/p/google-api-go-client/examples"
)

func not_init() {
	HttpClient := &http.Client{}
	//c := examples.getOAuthClient(config)
	driveService, err := drive.New(HttpClient)
	fileId := "1ne13EMak9-AEahjepgdN5z9scDkX21G5sYlOhcBaLaI"
	log.Println(driveService)
	f, err := driveService.Files.Get(fileId).Do()
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return
	}
	log.Printf("Title: %v", f.Title)
	log.Printf("Description: %v", f.Description)
	log.Printf("MIME type: %v", f.MimeType)
	log.Printf("Download Url: %v", f.DownloadUrl)
}

func googledriveservice_GetBlob(hash HCID) (b blob, err error) {
	return b, nil
}
func googledriveservice_GetCommit(hkid HKID) (c commit, err error) {
	return c, nil
}
func googledriveservice_GetKey(hkid []byte) (data blob, err error) {
	return data, nil
}
func googledriveservice_GetTag(hkid HKID, nameSegment string) (t Tag, err error) {
	return t, nil
}
