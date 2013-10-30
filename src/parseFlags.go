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

func takeActions(action *string, path *string) {
	if *action != "" {
		log.Println("action", *action)
	}
	if *path != "/" {
		log.Println("path", *path)
	}

	in := bufio.NewReader(os.Stdin)
	switch *action {
	case "createDomain":
		fmt.Println("Type in domain name:")
		domainName, _ := in.ReadString('\n')
		fmt.Printf("Name of Domain: %s", domainName)
	case "createRepository":
		fmt.Println("Type in repository name:")
		repositoryName, _ := in.ReadString('\n')
		fmt.Printf("Name of Repository: %s", repositoryName)
	case "insertDomain":
		fmt.Println("Type in domain name:")
		domainName, _ := in.ReadString('\n')
		fmt.Println("Insert HKID as a hexadecimal number:")
		hkid, _ := in.ReadString('\n')
		fmt.Printf("Name of Domain: %s", domainName)
		fmt.Printf("hkid: %s", hkid)
	case "insertRepository":
		fmt.Println("Type in repository name:")
		repositoryName, _ := in.ReadString('\n')
		fmt.Println("Insert HKID as a hexadecimal number:")
		hkid, _ := in.ReadString('\n')
		fmt.Printf("Name of Repository: %s", repositoryName)
		fmt.Printf("hkid: %s", hkid)
	}
}
