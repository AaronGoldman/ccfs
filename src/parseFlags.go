package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func parseFlagsAndTakeAction() {
	var mount = flag.Bool("mount", false, "Mount the fuse file system")
	var serve = flag.Bool("serve", true, "Start content object server")
	var createDomain = flag.Bool("createDomain", false, "Creates a new domain at path")
	var createRepository = flag.Bool("createRepository", false, "Creates a new repository at path")
	var insertDomain = flag.Bool("insertDomain", false, "Inserts the domain HKID at path")
	var insertRepository = flag.Bool("insertRepository", false, "Inserts the repository HKID at path")

	var path = flag.String("path", "", "The path to inserted collection")
	var hkid = flag.String("hkid", "", "HKID of collection to insert")

	flag.Parse()

	if flag.NFlag() == 0 {
		//flagged = false
		*serve = true
		*mount = true
	}

	if *serve {
		go BlobServerStart()
		go RepoServerStart()
	}
	if *mount {
		//go startFSintegration()
	}
	in := bufio.NewReader(os.Stdin)
	var err error
	h, collectionPath := fileSystemPath2CollectionPath(path)
	switch {
	case *createDomain:
		fmt.Println("Type in domain name:")
		domainName, _ := in.ReadString('\n')
		fmt.Printf("Name of Domain: %s", domainName)
		err = InitDomain(h, fmt.Sprintf("%s/%s", collectionPath, domainName))
		if err != nil {
			log.Println(err)
			return
		}

	case *createRepository:
		fmt.Println("Type in repository name:")
		repositoryName, _ := in.ReadString('\n')
		fmt.Printf("Name of Repository: %s", repositoryName)
		err = InitRepo(h, fmt.Sprintf("%s/%s", collectionPath, repositoryName))
		if err != nil {
			log.Println(err)
			return
		}
	case *insertDomain:
		fmt.Println("Type in domain name:")
		domainName, _ := in.ReadString('\n')
		fmt.Println("Insert HKID as a hexadecimal number:")
		var hex string = *hkid
		if *hkid == "" {
			hex, _ = in.ReadString('\n')
		}
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
	case *insertRepository:
		fmt.Println("Type in repository name:")
		repositoryName, _ := in.ReadString('\n')
		fmt.Println("Insert HKID as a hexadecimal number:")
		var hex string = *hkid
		if *hkid == "" {
			hex, _ = in.ReadString('\n')
		}
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

func fileSystemPath2CollectionPath(fileSystemPath *string) (h HKID, collectionPath string) {
	return HKID{}, ""
}
