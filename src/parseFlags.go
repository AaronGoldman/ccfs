package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func parseFlags() (*string, *string, bool) {
	var action = flag.String("action", "",
		"createDomain, createRepository, insertDomain, insertRepository")
	var path = flag.String("newPath", "/", "path where inserting happens")
	flag.Parse()
	flagged := true
	if flag.NFlag() == 0 {
		flagged = false
	}
	return action, path, flagged
}

func fileSystemPath2CollectionPath(fileSystemPath *string) (h HKID, collectionPath string) {

	return HKID{}, ""
}

func takeActions(action *string, fileSystemPath *string) {
	if *action != "" {
		log.Println("action", *action)
	}
	if *fileSystemPath != "/" {
		log.Println("path", *fileSystemPath)
	}

	in := bufio.NewReader(os.Stdin)
	var err error
	h, collectionPath := fileSystemPath2CollectionPath(fileSystemPath)
	switch *action {
	case "createDomain":
		fmt.Println("Type in domain name:")
		domainName, _ := in.ReadString('\n')
		fmt.Printf("Name of Domain: %s", domainName)
		err = InitDomain(h, fmt.Sprintf("%s/%s", collectionPath, domainName))
		if err != nil {
			log.Println(err)
			return
		}

	case "createRepository":
		fmt.Println("Type in repository name:")
		repositoryName, _ := in.ReadString('\n')
		fmt.Printf("Name of Repository: %s", repositoryName)
		err = InitRepo(h, fmt.Sprintf("%s/%s", collectionPath, repositoryName))
		if err != nil {
			log.Println(err)
			return
		}
	case "insertDomain":
		fmt.Println("Type in domain name:")
		domainName, _ := in.ReadString('\n')
		fmt.Println("Insert HKID as a hexadecimal number:")
		hex, _ := in.ReadString('\n')
		foreign_hkid, err := HkidFromHex(hex)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Name of Domain: %s", domainName)
		fmt.Printf("hkid: %s", h)
		err = InsertDomain(h, fmt.Sprintf("%s/%s", collectionPath, domainName), foreign_hkid)
		if err != nil {
			log.Println(err)
			return
		}
	case "insertRepository":
		fmt.Println("Type in repository name:")
		repositoryName, _ := in.ReadString('\n')
		fmt.Println("Insert HKID as a hexadecimal number:")
		hex, _ := in.ReadString('\n')
		foreign_hkid, err := HkidFromHex(hex)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Name of Repository: %s", repositoryName)
		fmt.Printf("hkid: %s", h)
		err = InsertRepo(h, fmt.Sprintf("%s/%s", collectionPath, repositoryName), foreign_hkid)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
