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

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

func parseFlags() (
	Flags struct {
		mount  *bool
		serve  *bool
		dht    *bool
		drive  *bool
		apps   *bool
		direct *bool
		lan    *bool
	},
	Command struct {
		createDomain     *bool
		createRepository *bool
		insertDomain     *bool
		insertRepository *bool
		path             *string
		hkid             *string
	},
) {

	Flags.mount = flag.Bool("mount", false, "Mount the fuse file system")
	Flags.serve = flag.Bool("serve", true, "Start content object server")
	Flags.dht = flag.Bool("dht", false, "Starts Kademliadht service")
	Flags.drive = flag.Bool("drive", false, "Starts Googledrive service")
	Flags.apps = flag.Bool("apps", false, "Starts Appsscript service")
	Flags.direct = flag.Bool("direct", false, "Starts direct http service")
	Flags.lan = flag.Bool("lan", false, "Starts the multicast service")
	Command.createDomain = flag.Bool("createDomain", false, "Creates a new domain at path argument")
	Command.createRepository = flag.Bool("createRepository", false, "Creates a new repository at path argument")
	Command.insertDomain = flag.Bool("insertDomain", false, "Inserts the domain HKID argument at path argument")
	Command.insertRepository = flag.Bool("insertRepository", false, "Inserts the repository HKID argument at path argument")

	Command.path = flag.String("path", "", "The path to inserted collection")
	Command.hkid = flag.String("hkid", "", "HKID of collection to insert")
	flag.Parse()
	return Flags, Command
}

func addCurators(inCommand struct {
	createDomain     *bool
	createRepository *bool
	insertDomain     *bool
	insertRepository *bool
	path             *string
	hkid             *string
}) {

	if *inCommand.path != "" {
		//log.Printf("HKID: %s", *hkid)
		in := bufio.NewReader(os.Stdin)
		var err error
		h, collectionPath := fileSystemPath2CollectionPath(inCommand.path)
		//log.Printf("systemPath %s", *path)
		//log.Printf("collectionPath %s", collectionPath)

		FileInfos, err := ioutil.ReadDir(*inCommand.path)
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

		collectionName := filepath.Base(*inCommand.path)
		//fmt.Printf("Name of Collection: %s\n", collectionName)

		switch {
		case *inCommand.createDomain:
			//err = InitDomain(h, fmt.Sprintf("%s/%s", collectionPath, collectionName))
			err = services.InitDomain(h, collectionPath)
			if err != nil {
				log.Println(err)
				return
			}

		case *inCommand.createRepository:
			//err = InitRepo(h, fmt.Sprintf("%s/%s", collectionPath, collectionName))
			err = services.InitRepo(h, collectionPath)
			if err != nil {
				log.Println(err)
				return
			}
		case *inCommand.insertDomain:
			fmt.Println("Insert HKID as a hexadecimal number:")
			hex := *inCommand.hkid
			if *inCommand.hkid == "" {
				hex, _ = in.ReadString('\n')
				hex = strings.Trim(hex, "\n")
			}
			log.Print(len(hex))
			foreignHkid, err := objects.HkidFromHex(hex)
			if err != nil {
				log.Printf("Somethng went wrong in insertDomain %s", err)
				os.Exit(2)
			}
			log.Printf("hkid: %s\n", h)
			err = services.InsertDomain(h, fmt.Sprintf("%s/%s", collectionPath, collectionName), foreignHkid)
			if err != nil {
				log.Println(err)
				return
			}
		case *inCommand.insertRepository:
			fmt.Println("Insert HKID as a hexadecimal number:")
			hex := *inCommand.hkid
			if *inCommand.hkid == "" {
				hex, _ = in.ReadString('\n')
				hex = strings.Trim(hex, "\n")
			}
			fmt.Printf("%s", hex)
			foreignHkid, err := objects.HkidFromHex(hex)
			if err != nil {
				log.Printf("Somethng went wrong in insertRepo %s", err)
				os.Exit(2)

			}
			fmt.Printf("hkid: %s", h)
			err = services.InsertRepo(h, fmt.Sprintf("%s/%s", collectionPath, collectionName), foreignHkid)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}

func fileSystemPath2CollectionPath(fileSystemPath *string) (objects.HKID, string) {
	h, err := objects.HkidFromHex("c09b2765c6fd4b999d47c82f9cdf7f4b659bf7c29487cc0b357b8fc92ac8ad02")
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
