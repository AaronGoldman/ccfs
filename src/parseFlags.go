package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	if *path != "" {
		log.Printf("HKID: %s", *hkid)
		in := bufio.NewReader(os.Stdin)
		var err error
		h, collectionPath := fileSystemPath2CollectionPath(path)
		log.Printf("systemPath %s", *path)
		log.Printf("collectionPath %s", collectionPath)

		FileInfos, err := ioutil.ReadDir(*path)
		if err != nil {
			log.Printf("Error reading directory %s", err)
			os.Exit(2)
			return
		}

		if len(FileInfos) != 0 {
			fmt.Printf("The folder is not empty")
			os.Exit(2)
			return // Ends function
		}

		collectionName := filepath.Base(*path)
		fmt.Printf("Name of Collection: %s\n", collectionName)

		switch {
		case *createDomain:
			err = InitDomain(h, fmt.Sprintf("%s/%s", collectionPath, collectionName))
			if err != nil {
				log.Println(err)
				return
			}

		case *createRepository:
			err = InitRepo(h, fmt.Sprintf("%s/%s", collectionPath, collectionName))
			if err != nil {
				log.Println(err)
				return
			}
		case *insertDomain:
			fmt.Println("Insert HKID as a hexadecimal number:")
			var hex string = *hkid
			if *hkid == "" {
				hex, _ = in.ReadString('\n')
			}
			foreign_hkid, err := HkidFromHex(hex)
			if err != nil {
				log.Printf("Somethng went wrong in insertDomain %s", err)
				os.Exit(2)
			}
			fmt.Printf("hkid: %s\n", h)
			err = InsertDomain(h, fmt.Sprintf("%s/%s", collectionPath, collectionName), foreign_hkid)
			if err != nil {
				log.Println(err)
				return
			}
		case *insertRepository:
			fmt.Println("Insert HKID as a hexadecimal number:")
			var hex string = *hkid
			if *hkid == "" {
				hex, _ = in.ReadString('\n')
			}
			fmt.Printf("%s", hex)
			foreign_hkid, err := HkidFromHex(hex)
			if err != nil {
				log.Printf("Somethng went wrong in insertRepo %s", err)
				os.Exit(2)

			}
			fmt.Printf("hkid: %s", h)
			err = InsertRepo(h, fmt.Sprintf("%s/%s", collectionPath, collectionName), foreign_hkid)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}

func fileSystemPath2CollectionPath(fileSystemPath *string) (HKID, string) {
	h, err := HkidFromHex("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02")
	wd, err := os.Getwd()
	path := filepath.Join(wd, "../mountpoint")
	//log.Printf("%s", *fileSystemPath)
	//LEFT OFF HERE- so fileSystemPath isnt what I want it to be
	*fileSystemPath = strings.Trim(*fileSystemPath, "\"")
	collectionPath, err := filepath.Rel(string(path), *fileSystemPath)
	if err != nil {
		log.Printf("OH NO! An Error %s", err)
		os.Exit(2)
	}

	return h, collectionPath
}
